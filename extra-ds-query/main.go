package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

// SampleModel is sample model.
type SampleModel struct {
	Name  string
	Value int
}

func main() {
	http.HandleFunc("/", login)

	// cloud_datastore.go
	http.HandleFunc("/store/get", getCloudDataStore)
	http.HandleFunc("/store/put", putCloudDataStore)

	// gae_datastore.go
	http.HandleFunc("/gaestore/get", getGAEDatastore)
	http.HandleFunc("/gaestore/put", putGAEDatastore)
	http.HandleFunc("/gaestore/not-equal", gaeNotEqual)
	http.HandleFunc("/gaestore/in", gaeIn)
	http.HandleFunc("/gaestore/in2", gaeIn2)
	http.HandleFunc("/gaestore/num-range", gaeNumRange)
	http.HandleFunc("/gaestore/like", gaeLike)
	http.HandleFunc("/gaestore/prefix", gaePrefix)

	// gcd_datastore.go
	http.HandleFunc("/gcdstore/not-equal", gcdNotEqual)
	http.HandleFunc("/gcdstore/in", gcdIn)
	http.HandleFunc("/gcdstore/in2", gcdIn2)
	http.HandleFunc("/gcdstore/num-range", gcdNumRange)
	http.HandleFunc("/gcdstore/like", gcdLike)
	http.HandleFunc("/gcdstore/prefix", gcdPrefix)

	// put testdata
	http.HandleFunc("/put-testbooks", putTestBooks)
	http.HandleFunc("/put-testgcdbooks", putTestGCDBooks) // ローカル実行用

	appengine.Main()
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if u == nil {
		url, _ := user.LoginURL(ctx, "/")
		fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
		return
	}
	url, _ := user.LogoutURL(ctx, "/")
	fmt.Fprintf(w, `Welcome, %#v! (<a href="%s">sign out</a>)`, u, url)
}
