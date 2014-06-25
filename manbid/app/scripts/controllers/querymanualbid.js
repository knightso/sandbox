'use strict';

angular.module('manbidApp')
  .controller('QuerymanualbidCtrl', function ($scope, Manualbid) {
    $scope.manualBids = Manualbid.query();
    $scope.criteria = {};
  });
