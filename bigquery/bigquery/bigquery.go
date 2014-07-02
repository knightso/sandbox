package bigquery

// OAuth2.0 Service Accounts ver.
/*
Google Client API Librariesのコードを一部変更しているのでメモ。
詳しくはブログに書いたので、忘れたら確認。

```go
// $GOPATH/src/code.google.com/p/google-api-go-client/bigquery/v2/bigquery-gen.go

type TableDataInsertAllRequestRows struct {
	InsertId string `json:"insertId,omitempty"`
	Json *JsonObject `json:"json,omitempty"`
}
```

これを次のように変更します。

```go
type TableDataInsertAllRequestRows struct {
	InsertId string `json:"insertId,omitempty"`
	Json interface{} `json:"json,omitempty"`    // こちらを変更。
}
```
*/

import (
	"appengine"
	"code.google.com/p/goauth2/appengine/serviceaccount"
	"code.google.com/p/google-api-go-client/bigquery/v2"
	"fmt"
	"net/http"
	"reflect"
)

const (
	PROJECT_ID string = "metal-bus-589"
	DATASET_ID string = "sample_dataset"
	TABLE_ID   string = "sample_table"
)

type Row struct {
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Population int    `json:"population"`
}

func init() {
	http.HandleFunc("/create_dataset", CreateDatasetHandler)
	http.HandleFunc("/create_table", CreateTableHandler)
	http.HandleFunc("/inset_all", InsertAllHandler)
	http.HandleFunc("/query", QueryHandler)
}

func CreateDatasetHandler(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	service, err := NewService(c)
	if err != nil {
		fmt.Fprintf(rw, "%s", err.Error())
		return
	}

	dataset, err := CreateDataset(service, PROJECT_ID, DATASET_ID)
	if err != nil {
		fmt.Fprintf(rw, "%s", err.Error())
		return
	}

	PrintStructField(rw, dataset)
}

func CreateTableHandler(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	service, err := NewService(c)
	if err != nil {
		fmt.Fprintf(rw, "%s", err.Error())
		return
	}

	schema := &bigquery.TableSchema{
		Fields: make([]*bigquery.TableFieldSchema, 3),
	}
	/*
		Typeに指定可能な値は、STRING, INTEGER, FLOAT, BOOLEAN, TIMESTAMP or RECORD。
		※RECORDはTableFieldSchema.Fieldsを使用している場合に使う。
	*/
	schema.Fields[0] = &bigquery.TableFieldSchema{
		Name: "kind",
		Type: "STRING",
	}
	schema.Fields[1] = &bigquery.TableFieldSchema{
		Name: "name",
		Type: "STRING",
	}
	schema.Fields[2] = &bigquery.TableFieldSchema{
		Name: "population",
		Type: "INTEGER",
	}
	table, err := CreateTable(service, PROJECT_ID, DATASET_ID, TABLE_ID, schema)
	if err != nil {
		fmt.Fprintf(rw, "%s", err.Error())
		return
	}

	PrintStructField(rw, table)
}

func InsertAllHandler(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	service, err := NewService(c)
	if err != nil {
		fmt.Fprintf(rw, "%s", err.Error())
		return
	}

	rows := make([]*Row, 2)
	rows[0] = &Row{Kind: "country", Name: "Shizuoka", Population: 707183}
	rows[1] = &Row{Kind: "country", Name: "Hamamatsu", Population: 791546}
	response, err := StreamingData(service, PROJECT_ID, DATASET_ID, TABLE_ID, rows)
	if err != nil {
		fmt.Fprintf(rw, "%s", err.Error())
		return
	}

	PrintStructField(rw, response)
}

func QueryHandler(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	service, err := NewService(c)
	if err != nil {
		fmt.Fprintf(rw, "%s", err.Error())
		return
	}

	query := "SELECT name FROM [" + DATASET_ID + "." + TABLE_ID + "] LIMIT 1000"
	response, err := JobQuery(service, PROJECT_ID, query)
	if err != nil {
		fmt.Fprintf(rw, "%s", err.Error())
		return
	}

	PrintStructField(rw, response)

	for _, row := range response.Rows {
		for _, cell := range row.F {
			fmt.Fprintf(rw, "%v", cell.V)
		}
	}
}

func PrintStructField(rw http.ResponseWriter, structure interface{}) {
	s := reflect.ValueOf(structure).Elem()
	typeOfStr := s.Type()
	for i := 0; i < s.NumField(); i++ {
		fmt.Fprintf(rw, "%s = %v\n", typeOfStr.Field(i).Name, s.Field(i).Interface())
	}
}

func NewService(c appengine.Context) (service *bigquery.Service, err error) {
	client, err := serviceaccount.NewClient(c, bigquery.BigqueryScope)
	if err != nil {
		return
	}

	service, err = bigquery.New(client)
	if err != nil {
		return
	}

	return service, nil
}

func CreateDataset(service *bigquery.Service, projectId, datasetId string) (dataset *bigquery.Dataset, err error) {
	dataset, err = service.Datasets.Insert(projectId, &bigquery.Dataset{
		DatasetReference: &bigquery.DatasetReference{
			DatasetId: datasetId,
			ProjectId: projectId,
		},
	}).Do()
	if err != nil {
		return
	}
	return dataset, nil
}

func CreateTable(service *bigquery.Service, projectId, datasetId, tableId string, schema *bigquery.TableSchema) (table *bigquery.Table, err error) {
	table, err = service.Tables.Insert(projectId, datasetId, &bigquery.Table{
		Schema: schema,
		TableReference: &bigquery.TableReference{
			DatasetId: datasetId,
			ProjectId: projectId,
			TableId:   tableId,
		},
	}).Do()
	if err != nil {
		return
	}
	return table, nil
}

func StreamingData(service *bigquery.Service, projectId, datasetId, tableId string, rows []*Row) (response *bigquery.TableDataInsertAllResponse, err error) {
	data := make([]*bigquery.TableDataInsertAllRequestRows, len(rows))
	for i, row := range rows {
		data[i] = &bigquery.TableDataInsertAllRequestRows{
			Json: row,
		}
	}
	response, err = service.Tabledata.InsertAll(projectId, datasetId, tableId, &bigquery.TableDataInsertAllRequest{
		Kind: "bigquery#tableDataInsertAllRequest",
		Rows: data,
	}).Do()
	if err != nil {
		return
	}
	return response, nil
}

func JobQuery(service *bigquery.Service, projectId, query string) (response *bigquery.QueryResponse, err error) {
	response, err = service.Jobs.Query(projectId, &bigquery.QueryRequest{
		Kind:  "bigquery#queryRequest",
		Query: query,
	}).Do()
	if err != nil {
		return
	}
	return response, nil
}
