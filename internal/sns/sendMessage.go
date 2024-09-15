package sns

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"go.uber.org/zap"
)

func getTopicArnFromName(sess *session.Session, topicName string) (*string, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	svc := sns.New(sess)

	// List SNS topics
	result, err := svc.ListTopics(nil)
	if err != nil {
		logger.Error("Error listing topics", zap.Error(err))
		return nil, err
	}

	logger.Info("Topics", zap.Any("topics", result.Topics))

	var topicArn *string
	for _, t := range result.Topics {
		if strings.HasSuffix(*t.TopicArn, ":"+topicName) {
			topicArn = t.TopicArn
			break
		}
	}

	return topicArn, nil
}

func SendMessage(topicName string, message interface{}) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String("eu-central-1")},
	}))

	svc := sns.New(sess)

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		logger.Error("Error marshalling message", zap.Error(err))
		return err
	}

	topicArn, err := getTopicArnFromName(sess, topicName)
	if err != nil {
		logger.Error("Error getting topic ARN", zap.Error(err))
		return err
	}

	_, err = svc.Publish(&sns.PublishInput{
		TopicArn: topicArn,
		Message:  aws.String(string(jsonMessage)),
	})

	if err != nil {
		logger.Error("Error sending message", zap.Error(err))
	}

	return err
}
