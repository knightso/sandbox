'use strict';
var gapi = require("googleapis");

var SCOPE = "https://spreadsheets.google.com/feeds";

var CLIENT_ID = '308896557542-kkgbsbionri51ktj1la3a0aqqir0csj4.apps.googleusercontent.com';
var CLIENT_SECRET = 'l26MgZWJPfhZ-7Ls7izN4w41';
var REDIRECT_URL = 'urn:ietf:wg:oauth:2.0:oob';

function Spreadsheet() {
	this.oAuth2Client = new gapi.OAuth2Client(CLIENT_ID, CLIENT_SECRET, REDIRECT_URL);
}

Spreadsheet.prototype.generateAuthUrl = function() {
	return this.oAuth2Client.generateAuthUrl({access_type: "offline", scope: SCOPE});
};

Spreadsheet.prototype.getAccessToken = function(code, callback) {
	var self = this;
	this.oAuth2Client.getToken(code, function(err, tokens) {
		self.oAuth2Client.setCredentials(tokens);
		callback(tokens);
	});
};

Spreadsheet.prototype.getWorksheets = function(spreadsheetKey, callback) {
	var opts = {
		url: SCOPE + "/worksheets/" + spreadsheetKey + "/private/full?alt=json"
	}
	this.oAuth2Client.request(opts, function(err, body, res) {
		if(err) {
			console.log(err);
			throw err;
		}
		callback(body);
	});
};

/*
Spreadsheet.prototype.getListFeed = function(worksheetId, callback) {
	var opts = {
		url: SCOPE + "/list/" + this.spreadsheetKey + "/" + worksheetId + "/private/basic?alt=json"
	}
	this.oAuth2Client.request(opts, function(err, body, res) {
		var result = {};
		result['entry'] = [];

		var entries = body.feed.entry;
		for(var i=0; i<entries.length; i++) {
			// 行データを分割整形する。
			var content = entries[i].content.$t;
			var contentArray = content.split(',');
			var row = {};
			for(var j=0; j<contentArray.length; j++) {
				var headerAndData = contentArray[j].split(':');
				var header = headerAndData[0];
				var cell = headerAndData[1];
				// 文字列先頭のスペースを削除し、セルの値を入れる。
				row[header.replace(/^\s+/g, '')] = cell.replace(/^\s+/g, '');
			}

			result['entry'].push({
				'title': entries[i].title.$t,
				'content': row
			});
		}
		callback(result);
	});
};
*/
Spreadsheet.prototype.getCells = function(spreadsheetKey, worksheetId, callback) {
	var opts = {
		url: SCOPE + "/cells/" + spreadsheetKey + "/" + worksheetId + "/private/full?alt=json"
	}
	this.oAuth2Client.request(opts, function(err, body, res) {
		callback(body);
	});
};

module.exports = Spreadsheet;