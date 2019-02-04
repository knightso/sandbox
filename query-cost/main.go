package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/taskqueue"

	"cloud.google.com/go/datastore"
)

func main() {
	http.HandleFunc("/store/query", queryHandler)
	http.HandleFunc("/store/keysOnly", keysOnlyHandler)
	http.HandleFunc("/store/projection", projectionHandler)

	http.HandleFunc("/store/addDataTask", putTestBooksTask)
	http.HandleFunc("/store/addData", putTestBooks)

	appengine.Main()
}

const (
	limit       = 1000000
	limitPerReq = 100000
	queueID     = "query-queue"
)

func queryHandler(w http.ResponseWriter, r *http.Request) {
	if !isTaskRequest(r) {
		if err := addTaskQueue(r, "/store/query", 0, ""); err != nil {
			log.Printf("error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		return
	}

	// get task parameter

	cursor := r.FormValue("cursor")
	count, err := strconv.Atoi(r.FormValue("count"))

	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// setup datastore

	ctx := r.Context()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// query

	q := datastore.NewQuery("Book2").Limit(limitPerReq)
	if c, err := datastore.DecodeCursor(cursor); err == nil {
		q = q.Start(c)
	}

	t := client.Run(ctx, q)
	for {
		var m Book2
		_, err := t.Next(&m)
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Printf("error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	c, err := t.Cursor()
	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	currentReadCount := count + limitPerReq
	log.Printf("query count: %d", currentReadCount)

	if c.String() != "" && currentReadCount < limit {
		if err = addTaskQueue(r, "/store/query", currentReadCount, c.String()); err != nil {
			log.Printf("error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func keysOnlyHandler(w http.ResponseWriter, r *http.Request) {
	if !isTaskRequest(r) {
		if err := addTaskQueue(r, "/store/keysOnly", 0, ""); err != nil {
			log.Printf("error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		return
	}

	// get task parameter

	cursor := r.FormValue("cursor")
	count, err := strconv.Atoi(r.FormValue("count"))

	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// setup datastore

	ctx := r.Context()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// query

	q := datastore.NewQuery("Book2").KeysOnly().Limit(limitPerReq)
	if c, err := datastore.DecodeCursor(cursor); err == nil {
		q = q.Start(c)
	}

	t := client.Run(ctx, q)
	for {
		_, err := t.Next(nil)
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Printf("error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	c, err := t.Cursor()
	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	currentReadCount := count + limitPerReq
	log.Printf("keysonly count: %d", currentReadCount)

	if c.String() != "" && currentReadCount < limit {
		if err = addTaskQueue(r, "/store/keysOnly", currentReadCount, c.String()); err != nil {
			log.Printf("error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func projectionHandler(w http.ResponseWriter, r *http.Request) {
	if !isTaskRequest(r) {
		if err := addTaskQueue(r, "/store/projection", 0, ""); err != nil {
			log.Printf("error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		return
	}

	// get task parameter

	cursor := r.FormValue("cursor")
	count, err := strconv.Atoi(r.FormValue("count"))

	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// setup datastore

	ctx := r.Context()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	client, err := datastore.NewClient(ctx, projectID)

	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// query

	q := datastore.NewQuery("Book2").Project("Title").Limit(limitPerReq)
	if c, err := datastore.DecodeCursor(cursor); err == nil {
		q = q.Start(c)
	}

	t := client.Run(ctx, q)
	for {
		_, err := t.Next(nil)
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Printf("error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	c, err := t.Cursor()
	if err != nil {
		log.Printf("error: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	currentReadCount := count + limitPerReq
	log.Printf("projection count: %d", currentReadCount)
	if c.String() != "" && currentReadCount < limit {
		if err = addTaskQueue(r, "/store/projection", currentReadCount, c.String()); err != nil {
			log.Printf("error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("done")
	}

	w.WriteHeader(http.StatusOK)
}

func addTaskQueue(r *http.Request, path string, count int, cursor string) error {
	ctx := appengine.NewContext(r)
	task := taskqueue.NewPOSTTask(path, url.Values{
		"count":  {strconv.Itoa(count)},
		"cursor": {cursor},
	})
	task.RetryOptions = &taskqueue.RetryOptions{RetryLimit: 0}
	_, err := taskqueue.Add(ctx, task, queueID)

	return err
}

func isTaskRequest(r *http.Request) bool {
	return r.Header.Get("X-AppEngine-QueueName") != ""
}
