(function() {
	var editedIDs = {};	// 編集したエディタのID一覧。

	tinymce.init({
		inline: true,
		language: 'ja',
		selector: 'div.editable',
		plugins: [
		  'advlist autolink lists link image charmap print preview anchor',
		  'searchreplace visualblocks code fullscreen',
		  'insertdatetime media table contextmenu paste catimgmanager',
		  'textcolor colorpicker'
		],
		toolbar: 'undo redo | styleselect | bold italic forecolor backcolor | alignleft aligncenter alignright alignjustify | bullist numlist outdent indent | link image catimgmanager',
		setup: function(editor) {
			// 編集したエディタを捕捉する。
			editor.on('blur', function(e) {
				console.log('blur');
				if(editor.id in editedIDs) {
					return;
				}
				editedIDs[editor.id] = null;
			});
		},
		file_picker_callback: function(callback, value, meta) {
			var editor = tinymce.activeEditor;
			// imageプラグインから参照ボタンを押した場合、ファイルブラウザを開きます。
			if(meta.filetype == 'image') {
				editor.windowManager.open({
					title: 'My File Browser',
					url: '/static/tinymce/plugins/catimgmanager/filebrowser/filebrowser.html',
					width: 800,
					height: 600,
					resizable: true,
					scrollbars: true,
					buttons: [
						{text: '閉じる', onclick: 'close'}
					]
				});
				editor.windowManager.setParams({callback: callback});
			}
		}
	});

	// エディタの編集したコンテンツをサーバーに送信します。
	window.onload = function() {
		var submitButton = document.getElementById('submitbutton');
		submitButton.addEventListener('click', function(evt) {
			var data = {'contents': {}};
			for(var id in editedIDs) {
				var editor = tinymce.EditorManager.get(id);
				// TODO: ブロック毎のIDを取得。
				data.contents[id] = editor.getContent();	// 今はエディタのIDだが、後々編集ブロック毎のIDにします。
			}

			tinymce.util.XHR.send({
				url: '/update/block',
				method: 'POST',
				data: data
			});
		}, false);
	}
})();