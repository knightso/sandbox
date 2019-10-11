package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
)

const (
	datasetID = "sample_dataset"
	tableID   = "sample_table"
)

type Row struct {
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Population int    `json:"population"`
}

func main() {
	http.HandleFunc("/create_dataset", CreateDatasetHandler)
	http.HandleFunc("/create_table", CreateTableHandler)
	http.HandleFunc("/insert_all", InsertAllHandler)
	http.HandleFunc("/query", QueryHandler)

	appengine.Main()
}

func CreateDatasetHandler(w http.ResponseWriter, r *http.Request) {
	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, proj)
	if err != nil {
		log.Fatal(err)
		return
	}

	meta := &bigquery.DatasetMetadata{
		Location: "asia-northeast1",
	}
	if err := client.Dataset(datasetID).Create(ctx, meta); err != nil {
		log.Fatal(err)
		return
	}
}

func CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, proj)
	if err != nil {
		log.Fatal(err)
		return
	}

	schema := bigquery.Schema{
		{Name: "kind", Type: bigquery.StringFieldType},
		{Name: "name", Type: bigquery.StringFieldType},
		{Name: "population", Type: bigquery.IntegerFieldType},
	}
	metaData := &bigquery.TableMetadata{
		Schema: schema,
	}

	tableRef := client.Dataset(datasetID).Table(tableID)
	if err := tableRef.Create(ctx, metaData); err != nil {
		log.Fatal(err)
		return
	}
}

func InsertAllHandler(w http.ResponseWriter, r *http.Request) {
	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, proj)
	if err != nil {
		log.Fatal(err)
		return
	}

	ins := client.Dataset(datasetID).Table(tableID).Inserter()
	rows := []Row{
		{Kind: "country", Name: "Shizuoka", Population: 707183},
		{Kind: "country", Name: "Hamamatsu", Population: 791546},
	}

	if err := ins.Put(ctx, rows); err != nil {
		log.Fatal(err)
		return
	}
}

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	proj := os.Getenv("GOOGLE_CLOUD_PROJECT")

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, proj)
	if err != nil {
		log.Fatal(err)
		return
	}

	query := client.Query("SELECT * FROM `" + datasetID + "." + tableID + "` LIMIT 1000")
	it, err := query.Read(ctx)

	if err != nil {
		log.Fatal(err)
		return
	}

	rows := []Row{}
	for {
		var row Row
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
			return
		}
		rows = append(rows, row)
	}

	b, err := json.Marshal(rows)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Write(b)
}
