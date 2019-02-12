package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/knightso/sandbox/tasktx"
	"google.golang.org/appengine"
	"google.golang.org/appengine/search"
)

func main() {
	http.HandleFunc("/putdata", handlePutData)

	appengine.Main()
}

func handlePutData(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("error = %s", err.Error())
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	var m tasktx.Model
	if err = json.Unmarshal(b, &m); err != nil {
		log.Printf("error = %s", err.Error())
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	index, err := search.Open("Model")
	if err != nil {
		log.Printf("error = %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	ctx := appengine.NewContext(r)
	_, err = index.Put(ctx, m.ID, &m)

	if err != nil {
		log.Printf("error = %s", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
