package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pietro-putelli/feynman-backend/internal/book"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"github.com/pietro-putelli/feynman-backend/internal/utility"
	"go.uber.org/zap"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	queryParams := &domain.BookListParams{}
	if err := utility.ParseQueryParams(request, queryParams); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		return *failure.NewBadRequest(fmt.Sprintf("Params validation failed: %s", err.Error())), nil
	}

	ctx, err := book.NewContext()
	if err != nil {
		logger.Error("Failed to create context", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	userID := utility.GetUserIDBy(request)

	switch queryParams.Type {
	case "short":
		books, err := ctx.Service.GetShortBooksList(userID, queryParams)
		if err != nil {
			logger.Error("Failed to get short books list", zap.Error(err))
			return *failure.NewInternalServerError(), nil
		}

		response, err := json.Marshal(books)
		if err != nil {
			logger.Error("Failed to marshal short books list response", zap.Error(err))
			return *failure.NewInternalServerError(), nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(response),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	case "long":
		books, err := ctx.Service.GetBooks(userID, queryParams)
		if err != nil {
			logger.Error("Failed to get books list", zap.Error(err))
			return *failure.NewInternalServerError(), nil
		}

		response, err := json.Marshal(books)
		if err != nil {
			logger.Error("Failed to marshal books list response", zap.Error(err))
			return *failure.NewInternalServerError(), nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(response),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       "Invalid type",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
