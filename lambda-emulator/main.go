package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"fmt"
)

var (
	sqsSvc   *sqs.SQS
	endpoint string = "http://sqs:9324"
	queueURL string = "http://sqs:9324/queue/events"
)

func pollMessages(chn chan<- *sqs.Message, plugin string) {

	for {
		output, err := sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: aws.Int64(2),
			WaitTimeSeconds:     aws.Int64(5),
		})

		if err != nil {
			fmt.Println("failed to fetch sqs message", err)
		}

		for _, message := range output.Messages {
			chn <- message
		}

	}

}

type Dag map[int][]Directive

type Directive struct {
	Plugin  string `json:"plugin"`
	Payload string `json:"payload"`
	Status  string `json:"status"`
}

func pluginID(messageBody *string) (string, error) {
	plugin := Directive{}
	if err := json.Unmarshal([]byte(*messageBody), &plugin); err != nil {
		return "", err
	}
	if plugin.Plugin == "" {
		return "", fmt.Errorf("missing or unknown `plugin` field in payload")
	}
	return plugin.Plugin, nil
}

// pullMessage...
func pullMessage(msg *sqs.Message) string {

	var dag Dag

	json.Unmarshal([]byte(string(*msg.Body)), &dag)

	for _, tasks := range dag {
		for _, task := range tasks {
			if task.Status == "complete" {
				continue
			} else {
				if task.Status == "blocked" {
					fmt.Println("Shipping job: ", task.Plugin)
					task.Status = "shipped"
					// send shipped job to process queue, blocked jobs to blocked
					return ""
				}
			}
			fmt.Println()
		}

	}

	pluginName, err := pluginID(msg.Body)
	if err != nil {
		return "unidentified plugin"
	}

	fmt.Println(pluginName, "message received", *msg.MessageId)

	return pluginName
}

func deleteMessage(msg *sqs.Message, pluginName string) {
	sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	fmt.Println(pluginName, "message deleted", *msg.MessageId)
}

func main() {

	plugins := []string{"hec-ras", "hydro-scalar", "consequences"}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		fmt.Println(err)
	}

	sqsSvc = sqs.New(sess, aws.NewConfig().WithEndpoint(endpoint))

	events := make(chan *sqs.Message, 2)
	for _, plugin := range plugins {
		go pollMessages(events, plugin)
	}

	for {
		for event := range events {
			pluginName := pullMessage(event)
			deleteMessage(event, pluginName)
		}

	}

}
