package server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Code copied and modified from:
// https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/go/sqs

func SendMsg(svc *sqs.SQS, queueURL *string, message string) (*sqs.SendMessageOutput, error) {

	sqsResponse, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(1),
		MessageBody:  aws.String(message),
		QueueUrl:     queueURL,
	})

	if err != nil {
		return sqsResponse, err
	}

	return sqsResponse, err
}
