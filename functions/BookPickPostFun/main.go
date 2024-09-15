package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
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

	body := &domain.CreateBookBody{}
	if err := json.Unmarshal([]byte(request.Body), body); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		return *failure.NewBadRequest("Invalid request body"), nil
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(body); err != nil {
		logger.Error("Body validation failed", zap.Error(err))
		return *failure.NewBadRequest("Body validation failed"), nil
	}

	ctx, err := book.NewContext()
	if err != nil {
		logger.Error("Failed to create context", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	bookData, err := ctx.Service.CreateBookPick(userID, body)
	if err != nil {
		logger.Error("Failed to create book or pick", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	response, err := json.Marshal(bookData)
	if err != nil {
		logger.Error("Failed to marshal response", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(response),
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
