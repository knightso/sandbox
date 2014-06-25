'use strict';

describe('Controller: QuerymanualbidCtrl', function () {

  // load the controller's module
  beforeEach(module('manbidApp'));

  var QuerymanualbidCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    QuerymanualbidCtrl = $controller('QuerymanualbidCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.awesomeThings.length).toBe(3);
  });
});
