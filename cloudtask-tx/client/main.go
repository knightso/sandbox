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

	rand.Seed(time.Now().UnixNano())
	val := rand.Float64()

	id := uuid.Must(uuid.NewV4()).String()

	log.Printf("putting id:%s\n", id)

	sample := tasktx.Sample{
		ID:        id,
		Title:     fmt.Sprintf("title %f", val),
		Value:     val,
		CreatedAt: time.Now(),
	}

	c := context.Background()

	client, err := datastore.NewClient(c, projectID)
	if err != nil {
		log.Fatal(err.Error())
	}

	c, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	_, err = client.RunInTransaction(c, func(tx *datastore.Transaction) error {

		key := datastore.NameKey("Sample", sample.ID, nil)
		if _, err := tx.Put(key, &sample); err != nil {
			log.Println("put model failed")
			return err
		}

		txStatus := &tasktx.TxStatus{
			ID:        uuid.Must(uuid.NewV4()).String(),
			CreatedAt: time.Now(),
		}

		txStatusKey := datastore.NameKey("TxStatus", txStatus.ID, nil)
		if _, err := tx.Put(txStatusKey, txStatus); err != nil {
			log.Println("put txStatus failed")
			return err
		}

		if err := startTask(c, txStatus.ID, sample); err != nil {
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

func startTask(ctx context.Context, txID string, sample tasktx.Sample) error {
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

	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2beta3#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2beta3#AppEngineHttpRequest
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
