package main

import (
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

// SampleModel is sample model.
type SampleModel struct {
	Name  string
	Value int
}

func getCloudDataStore(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	projectID := os.Getenv("STORE_PROJECT")
	client, err := datastore.NewClient(ctx, projectID, option.WithCredentialsFile("./serviceKey.json"))

	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	key := datastore.NameKey("SampleModel", "test", nil)

	var model SampleModel
	if err = client.Get(ctx, key, &model); err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fmt.Sprintf("%#v", model))
}

func putCloudDataStore(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	projectID := os.Getenv("STORE_PROJECT")
	client, err := datastore.NewClient(ctx, projectID, option.WithCredentialsFile("./serviceKey.json"))

	if err != nil {
		fmt.Fprintf(w, "error-01: %s", err.Error())
		return
	}

	key := datastore.NameKey("SampleModel", "test", nil)
	model := SampleModel{
		Name:  "Sample Model",
		Value: 123,
	}

	if _, err = client.Put(ctx, key, &model); err != nil {
		fmt.Fprintf(w, "error-02: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteCloudDataStore(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	projectID := os.Getenv("STORE_PROJECT")
	client, err := datastore.NewClient(ctx, projectID, option.WithCredentialsFile("./serviceKey.json"))

	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	key := datastore.NameKey("SampleModel", "test", nil)

	if err = client.Delete(ctx, key); err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func queryCloudDataStore(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	projectID := os.Getenv("STORE_PROJECT")
	client, err := datastore.NewClient(ctx, projectID, option.WithCredentialsFile("./serviceKey.json"))

	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
		return
	}
	q := datastore.NewQuery("SampleModel")
	t := client.Run(ctx, q)
	var list []SampleModel
	for {
		var m SampleModel
		_, err := t.Next(&m)
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Fprintf(w, "error: %s", err.Error())
			return
		}

		list = append(list, m)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fmt.Sprintf("%#v", list))
}
