'use strict';

angular.module('manbidApp')
  .controller('PutmanualbidCtrl', function ($scope, Manualbid, $routeParams) {
    if ($routeParams.sspId) {
      $scope.bid = Manualbid.get({
        sspId : $routeParams.sspId,
        mediaType : $routeParams.mediaType,
        mediaId : $routeParams.mediaId,
        campaignId : $routeParams.campaignId,
      });
    } else {
      $scope.bid = {
        sspId : $routeParams.sspId ? parseInt($routeParams.sspId, 10) : null,
        mediaType : $routeParams.mediaType ? $routeParams.mediaType : 'APP',
        mediaId : $routeParams.mediaId,
        campaignId : $routeParams.campaignId,
        cpm : null,
        isManualBid : true
      };
    }

    $scope.submit = function() {
      $scope.alerts = [];

      if (!$scope.updateManualBidForm.$valid) {
        return;
      }

      Manualbid.update($scope.bid,
        function(data) {
	  console.log('success: ' + angular.toJson(data, true));
          $scope.alerts.push({type: 'success', msg: 'successfully updated!'});
        },
        function(data) {
	  console.log('failure: ' + angular.toJson(data, true));
          $scope.alerts.push({type: 'danger', msg: 'update failed.'});
        }
      ); 
    };

    $scope.closeAlert = function(index) {
      $scope.alerts.splice(index, 1);
    }
  });
