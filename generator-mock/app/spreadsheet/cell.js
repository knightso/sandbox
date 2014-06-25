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
	var pattern = /^type=(string|int|float|date)$/i;
	var patternPk = /^(pk)\ntype=(string|int|float|date)$/;
	var result = {
		'primary': false,
		'type': null
	};

	var matched = content.match(pattern);
	if(matched != null) {
		result.type = matched[1];
	}else {
		matched = content.match(patternPk);
		if(matched != null) {
			result.primary = true;
			result.type = matched[2];
		}else {
			// TODO: 予期せぬ型情報としてエラー。
			console.log('content=[' + content + ']');
			throw new error.CellContentTypeException('有効な型ではありません。string, int, float, dateのいずれかにして下さい。', content);
		}
	}

	result.type = result.type.toLowerCase();
	return result;
};

module.exports = Cell;
