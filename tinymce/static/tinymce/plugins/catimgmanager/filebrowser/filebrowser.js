var app = angular.module('fileBrowserApp', [
	'fileBrowserControllers'
]);

var controller = angular.module('fileBrowserControllers', []);
controller.controller('imgListCtrl', ['$scope', '$http', function($scope, $http) {
	$scope.images = [];
	// 呼び出し元のTinyMCEのウィンドウからパラメータを取得します。
	var dialogArguments = window.top.tinymce.activeEditor.windowManager.getParams();
	var callback = dialogArguments['callback'];

	// GAEから画像一覧を取得。
	$http.get('/gcs/files')
		.success(function(data, status) {
			$scope.images = data['files'];
		});

	// 呼び出し元のTinyMCEのウィンドウに値を設定する。
	$scope.setDialog = function(index) {
		var img = $scope.images[index];
		callback(img.url, {alt: img.filename});
	};
}]);