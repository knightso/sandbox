'use strict';

function Error() {}

Error.prototype.TableAttributeConflictException = function(message, content) {
	this.message = message;
	this.content = content;
	this.name = 'TableAttributeConflictException';
};

Error.prototype.CellContentTypeException = function(message, content) {
	this.message = message;
  this.content = content;
  this.name = 'CellContentTypeException';
};

module.exports = Error;