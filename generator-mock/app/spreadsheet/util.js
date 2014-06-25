'use strict';
var fs = require('fs');
var path = require('path');

var spreadsheetError = require('./error.js');
var error = new spreadsheetError();

var TOKENS_FILEPATH = path.join(__dirname, 'tokens.json');

function Util() {}

Util.prototype.saveToken = function(tokens, callback) {
	fs.writeFile(TOKENS_FILEPATH, JSON.stringify(tokens), function(err) {
    if(err) {
    	callback(err);
    }else {
    	callback(null);
    }
  });
};

Util.prototype.readToken = function(callback) {
	fs.readFile(TOKENS_FILEPATH, function(err, data) {
		if(err) {
			callback(err, null);
		}else {
			callback(null, data);
		}
	});
};

Util.prototype.getWorksheetID = function(entryId) {
	var re = /^http(?:s)?:\/\/.+\/(.+)$/;
	var worksheetId = entryId.match(re)[1];
	return worksheetId;
};

Util.prototype.isAttributeConflict = function(attribute, cellContent) {
	try {
		for(var key in attribute) {
			if(cellContent == attribute[key]) {
				throw new error.TableAttributeConflictException('テーブルの属性が重複しています。', cellContent);
			}
		}
	}catch(err) {
		console.error(err.name, err.message, err.content);
		process.exit(1);
	}
};

Util.prototype.typeConversion = function(type, cellContent) {
	// cellContentをtypeで指定した型に変換する。
	var convertedContent;
	try {
    switch(type) {
      case 'string':
        convertedContent = cellContent;
        break;
      case 'int':
        if(isNaN(cellContent)) {
          throw new error.CellContentTypeException('int型に変換出来ない値です。', cellContent);;
        }else {
          convertedContent = parseInt(cellContent);
        }
        break;
      case 'float':
        if(isNaN(cellContent)) {
          throw new error.CellContentTypeException('float型に変換出来ない値です。', cellContent);;
        }else {
          convertedContent = parseFloat(cellContent);
        }
        break;
      case 'date':
        convertedContent = TransformDateFormat(cellContent);
        break;
      case null:
        convertedContent = cellContent;
        break;
    }
  }catch(err) {
    console.error(err.name, err.message, err.content);
    process.exit(1);
  }

  return convertedContent;
};

function TransformDateFormat(cellContent) {
	// spreadsheet apiでは、有効な日付書式ならばyyyy/MM/dd HHmmssで取得する。
	// 有効な場合、セルをクリックするとカレンダーが表示される。
	var longDatePattern = /^(\d{4})\/(\d{2})\/(\d{2})\s+(\d{1,2}):(\d{2}):(\d{2})/;
	var shortDatePattern = /^(\d{4})\/(\d{2})\/(\d{2})$/;
	var transformedDate = '';

	var dateArray = cellContent.match(longDatePattern);
	if(dateArray == null) {
		dateArray = cellContent.match(shortDatePattern);
	}
	
	if(dateArray == null) {
		throw new error.CellContentTypeException('有効な日付書式ではありません。Spreadsheetにてセル上をクリックした際、カレンダーが表示される書式に訂正して下さい。', content);
	}
	for(var i=1; i<dateArray.length; i++) {
		if(i == 4 && dateArray[i].length == 1) {
			// monthだけ1〜9月が一桁で取得されるので対処。
			transformedDate += '0';
			transformedDate += dateArray[i];
		}else {
			transformedDate += dateArray[i];
		}
	}

	return transformedDate;
}

module.exports = Util;
