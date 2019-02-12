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

const projectID = ""

const queueID = "put-data"

func main() {
	log.SetFlags(log.Lshortfile)

	rand.Seed(time.Now().UnixNano())
	val := rand.Float64()

	m := tasktx.Model{
		ID:    uuid.Must(uuid.NewV4()).String(),
		Title: fmt.Sprintf("title %f", val),
		Value: val,
	}

	ctx := context.Background()

	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		key := datastore.NameKey("Model", m.ID, nil)
		if _, err := tx.Put(key, &m); err != nil {
			return err
		}

		if err := startTask(ctx, m); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}

func startTask(ctx context.Context, m tasktx.Model) error {
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return err
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	queuePath := fmt.Sprintf("projects/%s/locations/us-central1/queues/%s", projectID, queueID)

	// Build the Task payload.
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
				},
			},
		},
	}

	_, err = client.CreateTask(ctx, req)
	return err
}
