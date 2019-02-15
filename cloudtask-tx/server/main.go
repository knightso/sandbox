package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/knightso/sandbox/tasktx"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"
)

func main() {
	http.HandleFunc("/putdata", handlePutData)

	appengine.Main()
}

func handlePutData(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	log.Printf("handlePutData. header=%v", r.Header)

	txID := r.Header.Get("X-TaskTx-ID")

	dispatchTime, err := time.Parse(time.RFC3339Nano, r.Header.Get("X-TaskTx-DispatchTime"))
	if err != nil {
		log.Printf("time.Parse failed:%s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	txStatusKey := datastore.NewKey(c, "TxStatus", txID, 0, nil)
	var txStatus tasktx.TxStatus
	if err = datastore.Get(c, txStatusKey, &txStatus); err == datastore.ErrNoSuchEntity {
		// トランザクション実行中またはタイムアウト
		if time.Now().Sub(dispatchTime) > 60*time.Second {
			log.Println("timeout")
			return
		} else {
			log.Println("retry")
			http.Error(w, err.Error(), http.StatusLocked)
			return
		}
	} else if err != nil {
		log.Printf("get TxStatus failed:%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Read Body failed:%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var sample tasktx.Sample
	if err = json.Unmarshal(b, &sample); err != nil {
		log.Printf("Unmarshal Body failed:%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get again just in case
	key := datastore.NewKey(c, "Sample", sample.ID, 0, nil)
	if err := datastore.Get(c, key, &sample); err != nil {
		log.Printf("get Sample failed:%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	index, err := search.Open("Sample")
	if err != nil {
		log.Printf("search.Open failed:%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = index.Put(c, sample.ID, &sample)

	if err != nil {
		log.Printf("put search index failed:%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("done")

	w.WriteHeader(http.StatusOK)
}
