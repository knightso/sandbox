ToolTip(Angular directive)
============

It is a tool tip using AngularJS directive.

----

次の２つのファイルを読み込んでください。

- tooltip.js
- tooltip.css

AngularJSのDIでtooltipDirectiveを渡してください。


```js
var app = angular.module('yourApp', [
	'tooltipDirective'
]);
```

ツールチップを出したい要素の直下に、&lt;tooltip&gt;を配置します。  
表示位置(tooltip-placement)には次の種類があります。

- top-center
- top-left
- top-right
- bottom-center
- bottom-left
- bottom-right

次の例では、&lt;button&gt;に対してツールチップを表示します。  
&lt;tooltip&gt;はbuttonの直下であればどこでも構いません。

```html
<button type="button">
	追加
	<tooltip tooltip-placement="top-right">単語を追加します</tooltip>
</button>
```