package main

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func getGAEDatastore(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	key := datastore.NewKey(ctx, "SampleModel", "gae-test", 0, nil)

	var model SampleModel
	if err := datastore.Get(ctx, key, &model); err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fmt.Sprintf("%#v", model))
}

func putGAEDatastore(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	key := datastore.NewKey(ctx, "SampleModel", "gae-test", 0, nil)
	model := SampleModel{
		Name:  "GAE Sample Model",
		Value: 987,
	}

	if _, err := datastore.Put(ctx, key, &model); err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
