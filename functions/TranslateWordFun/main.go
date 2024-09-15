package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"github.com/pietro-putelli/feynman-backend/internal/utility"
	"github.com/pietro-putelli/feynman-backend/langchain"
	"go.uber.org/zap"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	params := &domain.TranslateWordParams{}
	utility.ParseQueryParams(request, params)

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(params); err != nil {
		logger.Error("Params validation failed", zap.Error(err))
		return *failure.NewBadRequest("Params validation failed"), nil
	}

	response, err := langchain.TranslateWord(params.Word, params.Lang)
	if err != nil {
		logger.Error("Failed to translate word", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		logger.Error("Failed to marshal response body", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
