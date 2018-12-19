package main

import (
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
)

func getCloudDataStore(w http.ResponseWriter, r *http.Request) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := datastore.NewClient(r.Context(), projectID)

	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	key := datastore.NameKey("SampleModel", "test", nil)

	var model SampleModel
	if err = client.Get(r.Context(), key, &model); err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fmt.Sprintf("%#v", model))
}

func putCloudDataStore(w http.ResponseWriter, r *http.Request) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := datastore.NewClient(r.Context(), projectID)

	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	key := datastore.NameKey("SampleModel", "test", nil)
	model := SampleModel{
		Name:  "Sample Model",
		Value: 123,
	}

	if _, err = client.Put(r.Context(), key, &model); err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
