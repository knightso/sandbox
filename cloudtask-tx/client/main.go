package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	"cloud.google.com/go/datastore"
	"github.com/gofrs/uuid"
	"github.com/knightso/sandbox/tasktx"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"
)

const projectID = "metal-bus-589"

const queueID = "put-data"

func main() {
	log.SetFlags(log.Lshortfile)

	c := context.Background()

	client, err := datastore.NewClient(c, projectID)
	if err != nil {
		log.Fatal(err.Error())
	}

	sample := tasktx.Sample{
		ID:        uuid.Must(uuid.NewV4()).String(),
		Value:     rand.Float64(),
		CreatedAt: time.Now(),
	}

	// Taskのタイムアウトよりも短くする必要がある
	// Cloud Tasks APIは30sec以上のタイムアウトを指定できない
	c, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	_, err = client.RunInTransaction(c, func(tx *datastore.Transaction) error {

		// Sampleモデル保存
		key := datastore.NameKey("Sample", sample.ID, nil)
		if _, err := tx.Put(key, &sample); err != nil {
			log.Println("put model failed")
			return err
		}

		// Txステータスの保存
		txStatus := &tasktx.TxStatus{
			ID:        uuid.Must(uuid.NewV4()).String(),
			CreatedAt: time.Now(),
		}

		txStatusKey := datastore.NameKey("TxStatus", txStatus.ID, nil)
		if _, err := tx.Put(txStatusKey, txStatus); err != nil {
			log.Println("put txStatus failed")
			return err
		}

		// タスク起動
		if err := addTask(c, txStatus.ID, sample); err != nil {
			return err
		}

		//time.Sleep(40 * time.Second)

		return nil
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("done")
}

func addTask(ctx context.Context, txID string, sample tasktx.Sample) error {
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Println("cloudtasks NewClient failed")
		return err
	}

	b, err := json.Marshal(sample)
	if err != nil {
		return err
	}

	queuePath := fmt.Sprintf("projects/%s/locations/us-central1/queues/%s", projectID, queueID)

	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			PayloadType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
					HttpMethod:  taskspb.HttpMethod_POST,
					RelativeUri: "/putdata",
					Body:        b,
					Headers: map[string]string{
						"X-TaskTx-ID":           txID,
						"X-TaskTx-DispatchTime": time.Now().Format(time.RFC3339Nano),
					},
				},
			},
		},
	}

	_, err = client.CreateTask(ctx, req)
	if err != nil {
		log.Println("CreateTask failed")
		return err
	}

	return nil
}
