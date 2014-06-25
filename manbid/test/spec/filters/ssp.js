'use strict';

describe('Filter: ssp', function () {

  // load the filter's module
  beforeEach(module('manbidApp'));

  // initialize a new instance of the filter before each test
  var ssp;
  beforeEach(inject(function ($filter) {
    ssp = $filter('ssp');
  }));

  it('should return the input prefixed with "ssp filter:"', function () {
    var text = 'angularjs';
    expect(ssp(text)).toBe('ssp filter: ' + text);
  });

});
