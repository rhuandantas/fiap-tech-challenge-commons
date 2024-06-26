package messaging

import "context"

type Client interface {
	Publish(ctx context.Context, queueName string, message any) error
	Listen(ctx context.Context, queueUrl string, handler MessageHandler)
}

type MessageAttr map[string]string

type MessageHandler func(ctx context.Context, message string, attr MessageAttr) error
