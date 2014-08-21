'use strict';

var sperror = require('./error.js');
var error = new sperror();

function Cell() {};


Cell.prototype.getA1NotationPosition = function(title) {
	var pattern = /^([a-z]+)(\d+)$/i;	// セルアドレスを列・行に分割する。

	var notations = title.match(pattern);
	var column = notations[1];
	var row = parseInt(notations[2]);

	return {'column': column, 'row': row};
};

Cell.prototype.getColumnsType = function(content) {
	var pkPattern = /^pk\s*$/i;
	var typePattern = /^type=(\w+)\s*$/i;
	
	var result = {
		primary: false,
		type: 'string'
	};

	var lines = content.split('\n');

	for (var i = 0; i < lines.length; i++) {
		if (lines[i].match(pkPattern) != null) {
			result.primary = true
		}
		var matchedType = lines[i].match(typePattern);
		if (matchedType != null) {
			var type = matchedType[1];
			if (['string', 'number', 'date', 'json', 'eval'].indexOf(type) < 0) {
				console.log('content=[' + content + ']');
				throw new error.CellContentTypeException('有効な型ではありません。string, number, date, json, evalのいずれかにして下さい。', content);
			}
			result.type = type.toLowerCase();
		}
	}

	return result;
};

module.exports = Cell;
