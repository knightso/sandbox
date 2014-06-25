'use strict';

angular
  .module('manbidApp', [
    'ngResource',
    'ngRoute',
    'ui.bootstrap'
  ])
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .when('/putManualBid/:sspId;:mediaType;:mediaId;:campaignId', {
        templateUrl: 'views/putmanualbid.html',
        controller: 'PutmanualbidCtrl'
      })
      .when('/putManualBid', {
        templateUrl: 'views/putmanualbid.html',
        controller: 'PutmanualbidCtrl'
      })
      .when('/queryManualBid', {
        templateUrl: 'views/querymanualbid.html',
        controller: 'QuerymanualbidCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });
  });

var mocktangle = {}

mocktangle.importMockJs = function(jsfile) {
  if (location.hostname === 'mockhost') {
    document.write('<script type="text/javascript" src="' + jsfile + '"></script>');
  }
};

mocktangle.importMockJs('bower_components/angular-mocks/angular-mocks.js');
mocktangle.importMockJs('scripts/app-mock.js');

