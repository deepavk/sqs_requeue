package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"runtime/debug"
	"time"
)

var (
	destQueueUrl   *string
	sourceQueueUrl *string
	Svc            *sqs.SQS
	Region         *string
)

func getSession() *sqs.SQS {
	sess, _ := session.NewSession(&aws.Config{
		Region:     Region,
		MaxRetries: aws.Int(5),
	})
	Svc = sqs.New(sess)
	return Svc
}

func requeueMessages() error {
	conn := getSession()

	for {
		// receive max of 10 messages from dlq
		resp, err := conn.ReceiveMessage(&sqs.ReceiveMessageInput{
			WaitTimeSeconds:     aws.Int64(20),
			MaxNumberOfMessages: aws.Int64(10),
			QueueUrl:            sourceQueueUrl})
		if err != nil {
			log.Fatal(err)
			return err
		}
		messages := resp.Messages
		numberOfMessages := len(messages)
		log.Printf("Number of messsages to requeue %d", numberOfMessages)
		if numberOfMessages == 0 {
			return nil
		}

		// Send batch messages to destination queue
		log.Printf("Moving %v message(s) from %s to %s", numberOfMessages, *sourceQueueUrl, *destQueueUrl)
		var messageBatch []*sqs.SendMessageBatchRequestEntry
		for _, element := range messages {
			messageBatch = append(messageBatch,
				&sqs.SendMessageBatchRequestEntry{
					Id:                element.MessageId,
					MessageAttributes: element.MessageAttributes,
					MessageBody:       element.Body,
				})
		}
		log.Printf("Sending batch %+v to %s", messageBatch, *destQueueUrl)
		_, err = conn.SendMessageBatch(&sqs.SendMessageBatchInput{
			Entries:  messageBatch,
			QueueUrl: destQueueUrl})
		if err != nil {
			log.Fatal(err)
			return err
		}

		// delete batch messages from dlq
		var deleteMessageBatchRequestEntries []*sqs.DeleteMessageBatchRequestEntry
		for _, element := range messages {
			deleteMessageBatchRequestEntries = append(deleteMessageBatchRequestEntries,
				&sqs.DeleteMessageBatchRequestEntry{Id: element.MessageId,
					ReceiptHandle: element.ReceiptHandle})
		}
		log.Printf("Deleting batch %+v from %s", messageBatch, *sourceQueueUrl)
		_, err = conn.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
			Entries:  deleteMessageBatchRequestEntries,
			QueueUrl: sourceQueueUrl})
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
}

type RequeueEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name RequeueEvent) (string, error) {
	start := time.Now()
	err := requeueMessages()
	if err != nil {
		return "Error in requeue event", err
	}
	return fmt.Sprintf("Executed requeue event in %s", time.Since(start)), nil
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error in requeue: %s \n Stack trace: %s", err, string(debug.Stack()))
		}
	}()
	sourceQueueUrl = aws.String(os.Getenv("source_queue"))
	destQueueUrl = aws.String(os.Getenv("destination_queue"))
	Region = aws.String(os.Getenv("region"))
	lambda.Start(HandleRequest)
}
