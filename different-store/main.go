package main

import (
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	// cloud_datastore.go
	http.HandleFunc("/store/get", getCloudDataStore)
	http.HandleFunc("/store/put", putCloudDataStore)
	http.HandleFunc("/store/delete", deleteCloudDataStore)
	http.HandleFunc("/store/query", queryCloudDataStore)

	appengine.Main()
}
