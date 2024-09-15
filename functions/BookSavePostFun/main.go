package main

import (
	"encoding/json"

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

	userID := utility.GetUserIDBy(request)

	body := &domain.SaveBookBody{}

	if err := utility.ParseRequestBody(request, body); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		return *failure.NewBadRequest("Invalid request body"), nil
	}

	ctx, err := book.NewContext()
	if err != nil {
		logger.Error("Failed to create context", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	bookData, err := ctx.Service.SaveBook(userID, body)
	if err != nil {
		logger.Error("Failed to save book", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	response, err := json.Marshal(bookData)
	if err != nil {
		logger.Error("Failed to marshal response", zap.Error(err))
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
