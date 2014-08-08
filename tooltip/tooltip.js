
var tooltipDirective = angular.module('tooltipDirective', []);
tooltipDirective.directive('tooltip', [function() {
	function parsePlacement(placement) {
		result = placement.split('-');
		if(result.length != 2) {
			throw 'Invalid tooltip-placement.';
		}

		return {
			positionTB: result[0],
			positionLR: result[1]
		};
	}

	return {
		restrict: 'E',
		link: function(scope, iElement, iAttrs, controller) {
			// readyが無いと、cssが未適用の状態で幅や高さを取得してしまう。
			angular.element(document).ready(function() {
				// tooltipの親要素に対する処理
				var parent = iElement.parent();
				parent.css('position', 'relative');

				// クリックイベントの登録
				iElement.on("click", function() {
					iElement.css("display", "none");
				});

				// tooltip-placement属性の解析
				var placement = {};
				try {
					placement = parsePlacement(iAttrs['tooltipPlacement']);
				}catch(e) {
					console.log(e);
				}

				// tooltipのボックスオブジェクトを取得します。
				var tooltipRect = iElement[0].getBoundingClientRect();

				// tooltipの位置を左右に移動します。
				switch(placement.positionLR) {
					case 'right':
						var diff = parent[0].getBoundingClientRect().width - tooltipRect.width;
						iElement.css('left', diff);
						break;
					case 'left':
						// 何もしません。
						break;
					case 'center':
						var diff = parent[0].getBoundingClientRect().width - tooltipRect.width;
						iElement.css('left', diff / 2);
						break;
				}

				// tooltipの位置を上下に移動します。
				const space = 20;	// 親要素とtooltipの間隔。
				switch(placement.positionTB) {
					case 'top':
						iElement.css('top', -(tooltipRect.height + space));
						break;
					case 'bottom':
						iElement.css('top', (tooltipRect.height + space));
						break;
				}

				/* 矢印を追加します。 */
				iElement.prepend('<tooltiparrow class="tooltip-arrow"></tooltiparrow>');
				var arrow = iElement.find('tooltiparrow');
				var className = 'tooltip-arrow';
				var value = className;

				// 上下の配置を決めるクラスを追加します。
				switch(placement.positionTB) {
					case 'top':
						value = value + '-top';
						break;
					case 'bottom':
						value = value + '-bottom';
						break;
				}
				arrow.addClass(value);

				// x軸の配置を決めるクラスを追加します。
				switch(placement.positionLR) {
					case 'right':
						value = value + '-right';
						break;
					case 'left':
						value = value + '-left';
						break;
					case 'center':
						value = className + '-center';
						break;
				}
				arrow.addClass(value);
			});
		}
	}
}]);