package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
	"time"
)

const (
	AwsRegion = "AWS_REGION"
)

type SqsClient struct {
	client *sqs.Client
}

func NewSqsClient() SqsClient {
	//var (
	//	awsRegion = os.Getenv(AwsRegion)
	//)
	//
	//if awsRegion == "" {
	//	log.Fatalf("AWS_REGION environment variable not set")
	//}
	//sess := session.Must(session.NewSessionWithOptions(session.Options{
	//	SharedConfigState: session.SharedConfigEnable,
	//}))
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	//client := sqs.New
	client := sqs.NewFromConfig(cfg)

	return SqsClient{
		client: client,
	}
}

func (s SqsClient) Publish(ctx context.Context, queueUrl string, message any) error {
	messageStr, _ := json.Marshal(message)

	sendMessageInput := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueUrl),
		MessageBody: aws.String(string(messageStr)),
	}

	_, err := s.client.SendMessage(ctx, sendMessageInput)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to send message, %s", err.Error()))
	}

	return nil
}

func (s SqsClient) Listen(ctx context.Context, queueUrl string, handler MessageHandler) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Listener is shutting down...")
			return nil
		default:
			receiveMessageInput := &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queueUrl),
				MaxNumberOfMessages: 1,
				WaitTimeSeconds:     10, // Long polling
			}

			result, err := s.client.ReceiveMessage(ctx, receiveMessageInput)
			if err != nil {
				log.Printf("Failed to receive messages: %v", err)
				time.Sleep(5 * time.Second) // Delay before retrying
				continue
			}

			if len(result.Messages) == 0 {
				fmt.Println("No messages received")
				time.Sleep(2 * time.Second) // Delay before next receive attempt
				continue
			}

			for _, message := range result.Messages {
				fmt.Printf("Message received: %s\n", *message.Body)

				err = handler(ctx, *message.Body, nil)
				if err != nil {
					log.Printf("Failed to handle message: %v", err)
					time.Sleep(5 * time.Second)
					return nil
				}

				deleteMessageInput := &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(queueUrl),
					ReceiptHandle: message.ReceiptHandle,
				}

				_, err = s.client.DeleteMessage(ctx, deleteMessageInput)
				if err != nil {
					log.Printf("Failed to delete message: %v", err)
					// Optionally, add message to a dead-letter queue or retry deleting
				} else {
					fmt.Println("Message deleted successfully")
				}
			}
		}
	}
}
