'use strict';

describe('Service: Manualbid', function () {

  // load the service's module
  beforeEach(module('manbidApp'));

  // instantiate service
  var Manualbid;
  beforeEach(inject(function (_Manualbid_) {
    Manualbid = _Manualbid_;
  }));

  it('should do something', function () {
    expect(!!Manualbid).toBe(true);
  });

});
