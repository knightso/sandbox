package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"github.com/gofrs/uuid"
)

const reqTopic = "grush-req"

// Request represents test request.
type Request struct {
	ID   string
	Args []string
}

// TestResult represents a test result.
type TestResult struct {
	ID        string
	Output    string `datastore:",noindex"`
	Error     string `datastore:",noindex"`
	CreatedAt time.Time
}

// NewTestResultKey returns TestResult key
func NewTestResultKey(id string) *datastore.Key {
	return datastore.NameKey("TestResult", id, nil)
}

var pbClient *pubsub.Client
var dsClient *datastore.Client

var project string
var targetURL string

func init() {
	project = os.Getenv("GRUSH_GCP_PROJECT")
	if project == "" {
		log.Fatalf("need to set environment variable:GRUSH_GCP_PROJECT\n")
	}

	targetURL = os.Getenv("GRUSH_TARGET_URL")
	if targetURL == "" {
		log.Fatalf("need to set environment variable:GRUSH_TARGET_URL\n")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, project)
	if err != nil {
		log.Fatalf("failed to create pubsub client. err:%v", err)
	} else {
		pbClient = client
	}

	if client, err := datastore.NewClient(ctx, project); err != nil {
		log.Fatalf("failed to create datastore client. err:%v", err)
	} else {
		dsClient = client
	}
}

func main() {
	ctx := context.Background()

	topic := createTopicIfNotExists(pbClient, reqTopic)

	// subscription名 は英字から始まる必要があるので a- を先頭につける
	subID := "a-" + uuid.Must(uuid.NewV4()).String()
	sub, err := pbClient.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		log.Fatalf("failed to create subscription. err:%v", err)
	}
	defer sub.Delete(ctx)

	fmt.Printf("create subscription: %s\n", subID)
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()

		fmt.Printf("request found: %+v\n", msg)

		var req Request
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			log.Printf("could not decode message data: %#v", msg)
			return
		}

		output, err := runTest(&req)

		fmt.Printf("test output: %s\n", output)

		result := TestResult{
			ID:        req.ID,
			Output:    output,
			CreatedAt: time.Now(),
		}

		if err != nil {
			result.Error = err.Error()
		}

		key := NewTestResultKey(req.ID)
		_, err = dsClient.Put(ctx, key, &result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "save result failed: %s\n", err.Error())
		}

		// cancel をしないと次の Publish を待ち続けるので cancel して終了する。
		cancel()
	})

	if err != nil {
		log.Fatalf("failt to subscribe. err: %v", err)
	}
}

func createTopicIfNotExists(c *pubsub.Client, topicName string) *pubsub.Topic {
	ctx := context.Background()

	// Create a topic to subscribe to.
	t := c.Topic(topicName)
	ok, err := t.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		return t
	}

	t, err = c.CreateTopic(ctx, topicName)
	if err != nil {
		log.Fatalf("Failed to create the topic: %v", err)
	}
	return t
}

func runTest(req *Request) (output string, err error) {
	args := append(req.Args, targetURL)

	fmt.Printf("ARGS: %#v", args)
	cmd := exec.Command("ab", args...)

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Run(); err != nil {
		return buf.String(), err
	}

	return buf.String(), nil
}
