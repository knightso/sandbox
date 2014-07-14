/*
logパッケージの動作確認用httpハンドラを定義したファイルです。
*/

package logging

import (
	//"appengine"
	"fmt"
	"net/http"
)

const (
	datasetName string = "test"
	queueName   string = "logQueue"
)

func init() {
	// スキーマの設定。
	AddSchema("table_a", "kind:string,date:timestamp,count:integer")
	AddSchema("table_b", "kind:string,date:timestamp,count:integer")

	http.HandleFunc("/schema_check", SchemaCheckHandler)
	//http.HandleFunc("/put_log", PutLogHandler)
	//http.HandleFunc("/schedule", ScheduleHandler)
}

func SchemaCheckHandler(rw http.ResponseWriter, req *http.Request) {
	for table, schema := range schemata {
		fmt.Fprintf(rw, "Table name: %s\n", table)
		for column, typeName := range schema {
			fmt.Fprintf(rw, "%s: %s\n", column, typeName)
		}
		fmt.Fprintf(rw, "\n")
	}
}

/*
func PutLogHandler(rw http.ResponseWriter, req *http.Request) {
	c := appengine.NewContext(req)
	service := &LogService{Context: c}
	err := service.Put(queueName, tag, record)
}

func ScheduleHandler(rw http.ResponseWriter, req *http.Request) {
	// taskqueueに溜まったtaskを取り出して、BigQueryに格納します。
}
*/
