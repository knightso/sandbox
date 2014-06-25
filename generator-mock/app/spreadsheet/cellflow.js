'use strict';
var async = require('async');

var spreadsheetCell = require('./cell.js');
var spreadsheetUtil = require('./util.js');

var cell = new spreadsheetCell();
var ssutil = new spreadsheetUtil();

function CellFlow(client) {
	this.client = client;
	this.spreadsheetKey = '';
}

CellFlow.prototype.main = function(spreadsheetKey, ev) {
	this.spreadsheetKey = spreadsheetKey;
  var self = this;

	this.client.getWorksheets(this.spreadsheetKey, function(body) {
    self.worksheet(body, ev);
  });

  ev.once('cellflow_complete', function(data) {
  	ev.emit('return_data', data);
  });
};

CellFlow.prototype.worksheet = function(body, ev) {
	var self = this;
	var spreadsheet = {
    'title': body.feed.title.$t,
    'tables': {}
  };
  var worksheets = body.feed.entry;

  async.each(worksheets, function(worksheet, callback) {
    var worksheetTitle = worksheet.title.$t;
    var worksheetId = ssutil.getWorksheetID(worksheet.id.$t);
    self.client.getCells(self.spreadsheetKey, worksheetId, function(body) {
      spreadsheet['tables'][worksheetTitle] = self.cell(body);
      callback(null);
    });
  }, function(err) {
    if(err) {
      console.log(err);
      process.exit(1);
    }
    ev.emit('cellflow_complete', spreadsheet);
  });
};

CellFlow.prototype.cell = function(body) {
	var worksheet = {'records': [], 'primary_key': []};
  var attr = {};
  var type = {};

  var cells = body.feed.entry;
  // 属性と主キーの取得。
  for(var i=0; i<cells.length; i++) {
    var content = cells[i].content.$t;
    var position = cell.getA1NotationPosition(cells[i].title.$t);

    if(position.row == 1) {
      // 属性の重複を確認し、新しい列の属性を記録する。
      ssutil.isAttributeConflict(attr, content);
      attr[position.column] = content;
    }else if(position.row == 2) {
      // 列の型を記録し、Primary Key列であれば別途記録する。
      var columnsType = cell.getColumnsType(content);
      if(columnsType.primary) {
        worksheet['primary_key'].push(attr[position.column]);
      }
      type[position.column] = columnsType.type;
    }
  }

  // 型情報が空のものをデフォルトのString型に設定する。
  var attrbuteKeys = Object.keys(attr);
  for (var i = 0; i < attrbuteKeys.length; i++) {
    var column = attrbuteKeys[i];
    if(column in type) {
      // pass
    }else {
      type[column] = 'string';
    }
  }

  // テーブルを作成する。
  var records = [];
  var row = {};
  var rowNum = 3;
  for (var i = 0; i < cells.length; i++) {
    var position = cell.getA1NotationPosition(cells[i].title.$t);
    if(position.row == 1 || position.row == 2) {
      continue;
    }
    // 改行したので、前の行を記録してからrowを消す。
    if(rowNum != position.row) {
      records.push(row);
      row = {};
      rowNum = position.row;
    }
    var content = cells[i].content.$t;
    row[attr[position.column]] = ssutil.typeConversion(type[position.column], content);
  }
  records.push(row);

  worksheet['records'] = records;
  return worksheet;
};

module.exports = CellFlow;