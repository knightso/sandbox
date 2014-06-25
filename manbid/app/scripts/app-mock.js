'use strict';

(function () {
  $('body[ng-app]').attr('ng-app', function(arr) {
    return $(this).attr('ng-app') + 'Dev';
  });
  
  // モックappを作成して登録（名前は対象app名 + 'Dev'）
  var appDev = angular.module('manbidAppDev', ['manbidApp', 'ngMockE2E'/*, 'mockCommon'*/]);
  appDev.run(['$httpBackend', function($httpBackend) {
    
    // 正規表現にマッチするパスに対してXHRがあったらモックレスポンス(JSON)を返却
    $httpBackend.whenPUT(/\/api\/manualbid\/.*/).respond({result: 'success'});

    // ↓ $resourceでslashで終わるURLを指定する場合はbackslashでescapeするのだが、
    // $httpBackendで指定する場合には下記の様にしないとhookできない。bug?
    $httpBackend.whenGET('/api/manualbid\\').respond([
      {
        sspId : 19,
        mediaType : 'APP',
        mediaId : 'media123',
        campaignId : 'cmp000001',
        cpm : 123.456,
        isManualBid : true
      },
      {
        sspId : 20,
        mediaType : 'SITE',
        mediaId : 'media456',
        campaignId : 'cmp000002',
        cpm : 100,
        isManualBid : true
      },
      {
        sspId : 21,
        mediaType : 'APP',
        mediaId : 'media123',
        campaignId : 'cmp000003',
        cpm : 110.999,
        isManualBid : true
      },
      {
        sspId : 21,
        mediaType : 'SITE',
        mediaId : 'media456',
        campaignId : 'cmp000004',
        cpm : 120,
        isManualBid : true
      },
      {
        sspId : 19,
        mediaType : 'APP',
        mediaId : 'media123',
        campaignId : 'cmp000005',
        cpm : 123.456,
        isManualBid : false
      },
      {
        sspId : 20,
        mediaType : 'SITE',
        mediaId : 'media456',
        campaignId : 'cmp000006',
        cpm : 100,
        isManualBid : true
      },
      {
        sspId : 21,
        mediaType : 'APP',
        mediaId : 'media123',
        campaignId : 'cmp000007',
        cpm : 110.999,
        isManualBid : true
      },
      {
        sspId : 21,
        mediaType : 'SITE',
        mediaId : 'media456',
        campaignId : 'cmp000008',
        cpm : 120,
        isManualBid : false
      }
    ]);
   
    $httpBackend.whenGET(/\/manualbid\/.+/).respond({
      sspId : 21,
      mediaType : 'APP',
      mediaId : 'media123',
      campaignId : 'cmp000007',
      cpm : 110.999,
      isManualBid : true
    });

    // htmlファイルの取得等はそのままスルー
    $httpBackend.whenGET(/views\/.*/).passThrough();
  }]);
}());

