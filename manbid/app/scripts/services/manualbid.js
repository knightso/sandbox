'use strict';

angular.module('manbidApp')
  .service('Manualbid', function Manualbid($resource) {
    return $resource('/api/manualbid/:sspId/:mediaType/:mediaId/:campaignId', 
        {sspId:'@sspId', mediaType:'@mediaType', mediaId:'@mediaId', campaignId:'@campaignId'}, {
      query: {method:'GET', url:'/api/manualbid\\/', isArray:true},
      update: {method:'PUT'}
    });
  });
