package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pietro-putelli/feynman-backend/internal/book"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"github.com/pietro-putelli/feynman-backend/internal/utility"
	"go.uber.org/zap"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	userID := utility.GetUserIDBy(request)

	ctx, err := book.NewContext()
	if err != nil {
		logger.Error("Failed to create context", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	topics, err := ctx.Service.GetUserBooksTopics(userID)
	if err != nil {
		logger.Error("Failed to get topics", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	response, err := json.Marshal(topics)
	if err != nil {
		logger.Error("Failed to marshal topics", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(response),
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
