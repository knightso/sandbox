package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/gofrs/uuid"
)

const reqTopic = "grush-req"

// Request represents test request.
type Request struct {
	ID   string
	Args []string
}

var project string

var pbClient *pubsub.Client

func init() {
	project = os.Getenv("GRUSH_GCP_PROJECT")
	if project == "" {
		log.Fatalf("need to set environment variable:GRUSH_GCP_PROJECT\n")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("failed to create pubsub client. err:%v", err)
	}

	pbClient = client
}

func main() {
	id := uuid.Must(uuid.NewV4()).String()
	args := os.Args[1:len(os.Args)]
	if err := doRequest(id, args); err != nil {
		log.Fatalf("request failed: %v", err)
	}
}

func doRequest(id string, args []string) error {
	req := Request{
		ID:   id,
		Args: args,
	}
	return insertRequest(&req)
}

func insertRequest(req *Request) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	ctx := context.Background()

	topic := pbClient.Topic(reqTopic)
	_, err = topic.Publish(ctx, &pubsub.Message{Data: data}).Get(ctx)

	if err != nil {
		return err
	}

	return nil
}
