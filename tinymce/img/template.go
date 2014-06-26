package img

import (
	"appengine"
	"encoding/json"
	"net/http"
)

type TemplateMeta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Url         string `json:"url"`
}

func init() {
	http.HandleFunc("/tinymce/templatelist", TemplateListHandler)
}

// TinyMCEのtemplateプラグインに対して、テンプレート情報の一覧を返します。
func TemplateListHandler(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	metas, err := GetTemplateMeta()
	if err != nil {
		c.Errorf("%s", err.Error())
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Write(metas)
}

func GetTemplateMeta() (metas []byte, err error) {
	dammy1 := &TemplateMeta{
		Title:       "Greeting",
		Description: "Greeting message",
		Content:     "Hello world!!",
	}
	dammy2 := &TemplateMeta{
		Title:       "Form",
		Description: "File form",
		Url:         "/static/tmpl/fileform.html",
	}

	metas, err = json.Marshal([]*TemplateMeta{dammy1, dammy2})
	if err != nil {
		return []byte{}, err
	}
	return
}
