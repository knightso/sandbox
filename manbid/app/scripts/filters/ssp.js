'use strict';

angular.module('manbidApp')
  .filter('ssp', function (Enums) {
    return function (sspId) {
      var ssp = Enums.findSsp(sspId);
      console.log('ssp='+ssp);
      return ssp ? ssp.name + '(' + ssp.id + ')' : 'err';
    };
  });
