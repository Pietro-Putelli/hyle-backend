package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pietro-putelli/feynman-backend/internal/auth"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"go.uber.org/zap"
)

func handler(handleFunc auth.AuthTokenHandler) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		body := &domain.AuthTokenBody{}
		if err := json.Unmarshal([]byte(request.Body), body); err != nil {
			logger.Error("Invalid request body", zap.Error(err))
			return *failure.NewBadRequest("Invalid request body"), nil
		}

		context, err := auth.NewBaseContext()
		if err != nil {
			logger.Error("Failed to create context", zap.Error(err))
			return *failure.NewInternalServerError(), nil
		}

		tokenInfo, err := handleFunc(context, body)
		if err != nil {
			// handle validator error
			if errors.Is(err, &failure.ValidationErr{}) {
				logger.Error("Validation error", zap.Error(err))
				return *failure.NewBadRequest(err.Error()), nil
			}

			// handle provider error
			if errors.Is(err, &auth.InvalidProviderError{}) {
				logger.Error("Invalid provider", zap.Error(err))
				return *failure.NewBadRequest("Invalid provider"), nil
			}

			// handle internal server error
			logger.Error("Failed to handle auth token", zap.Error(err))
			return *failure.NewInternalServerError(), nil
		}

		response, err := json.Marshal(tokenInfo)
		if err != nil {
			logger.Error("Failed to marshal response", zap.Error(err))
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
}

func main() {
	lambda.Start(handler(auth.HandleAuthToken))
}
