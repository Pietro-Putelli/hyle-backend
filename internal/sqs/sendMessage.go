package sqs

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func getQueueURL(sess *session.Session, queue *string) (*sqs.GetQueueUrlOutput, error) {
	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func SendMessage(queueName string, message interface{}) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String("eu-central-1")},
	}))

	svc := sqs.New(sess)

	queueURL, err := getQueueURL(sess, &queueName)
	if err != nil {
		return err
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    queueURL.QueueUrl,
		MessageBody: aws.String(string(jsonMessage)),
	})

	return err
}
