package another

// OAuth2.0 Installed Apps ver.
/*
こちらは、Google Developer ConsoleでAPI & AUTH -> CredentialsページでInstalled Apps用のクライアントID
を作成して、OAuth2.0認証を行う場合のコードです。
*/

import (
	"appengine"
	"appengine/urlfetch"
	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/bigquery/v2"
	"fmt"
	"net/http"
)

var config = &oauth.Config{
	ClientId:     "961460936936-nr7g2ssks78k06c25k2o94lrfu9q2u0f.apps.googleusercontent.com",
	ClientSecret: "LV6GlMgwEb4gSqGB3UdLMkKU",
	Scope:        bigquery.BigqueryScope,
	RedirectURL:  "http://localhost:8080/bigquery",
	AuthURL:      "https://accounts.google.com/o/oauth2/auth",
	TokenURL:     "https://accounts.google.com/o/oauth2/token",
}

type People struct {
	Kind     string `json:"kind"`
	FullName string `json:"fullName"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
}

func init() {
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/bigquery", BQHandler)
}

func RootHandler(rw http.ResponseWriter, req *http.Request) {
	url := config.AuthCodeURL("")
	http.Redirect(rw, req, url, http.StatusFound)
}

func BQHandler(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	code := req.FormValue("code")
	t := &urlfetch.Transport{Context: c}
	transport := &oauth.Transport{
		Config:    config,
		Transport: t,
	}
	_, err := transport.Exchange(code)
	if err != nil {
		fmt.Fprintf(rw, "%s", err.Error())
		return
	}
	client := transport.Client()
	service, err := bigquery.New(client)
	if err != nil {
		fmt.Fprintf(rw, "%s", err.Error())
		return
	}

	rows := make([]*bigquery.TableDataInsertAllRequestRows, 1)
	/*
		rows[0] = &bigquery.TableDataInsertAllRequestRows{
			Json: &bigquery.JsonObject{
				kind:     "person",
				fullName: "Eddy Kingston",
				age:      26,
				gender:   "Male",
			},
		}
	*/
	rows[0] = &bigquery.TableDataInsertAllRequestRows{
		Json: People{
			Kind:     "person",
			FullName: "Eddy Kingston",
			Age:      26,
			Gender:   "Male",
		},
	}

	data := &bigquery.TableDataInsertAllRequest{
		Kind: "bigquery#tableDataInsertAllRequest",
		Rows: rows,
	}
	_, err = service.Tabledata.InsertAll(
		"metal-bus-589", // projectId
		"test_data",     // datasetId
		"test_table",    // tableId
		data,            // tabledatainsertallrequest
	).Do()
	if err != nil {
		fmt.Fprintf(rw, "%s", err.Error())
		return
	}
	fmt.Printf("%s", "Success.")
}
