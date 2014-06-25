'use strict';

angular.module('manbidApp')
  .service('Enums', function Enums() {
    return {
      ssps : [
        {name:'smaato', id:1},
        {name:'mopub', id:2},
        {name:'nexage', id:3},
        {name:'mobclix', id:4},
        {name:'pubmatic', id:5},
        {name:'amobee', id:6},
        {name:'tapit', id:9},
        {name:'millennial', id:11},
        {name:'flurry', id:14},
        {name:'geniee', id:15},
        {name:'ninja', id:16},
        {name:'adstir', id:18},
        {name:'mediba', id:19},
        {name:'gmo', id:20},
        {name:'glossom', id:21}
      ],
      findSsp : function(id) {
        for (var i = 0; i < this.ssps.length; i++) {
          if (this.ssps[i].id === id) return this.ssps[i];
        }
        return null;
      }
    };
  })
  .run(function ($rootScope, Enums) {
    $rootScope.Enums = Enums;
  });
