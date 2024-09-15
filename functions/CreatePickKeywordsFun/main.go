package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pietro-putelli/feynman-backend/internal/book"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/langchain"
	"go.uber.org/zap"
)

func handler(context context.Context, sqsEvent events.SQSEvent) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	messageBody := sqsEvent.Records[0].Body
	message := &domain.BookPickSearchKeywordMessage{}

	err := json.Unmarshal([]byte(messageBody), &message)
	if err != nil {
		log.Fatalf("Failed to unmarshal SQS message: %v", err)
	}

	ctx, err := book.NewContext()
	if err != nil {
		logger.Error("Error creating new context", zap.Error(err))
		return err
	}

	// keywords := []string{}

	keywords, err := langchain.GeneratePickKeywords(message.PickContent)
	if err != nil {
		logger.Error("Error generating pick keywords", zap.Error(err))
		return err
	}

	err = ctx.Service.AddPickKeywords(message.UserGuid, message.PickID, keywords)
	if err != nil {
		logger.Error("Error adding pick keywords", zap.Error(err))
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
