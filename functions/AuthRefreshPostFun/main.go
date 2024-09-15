package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pietro-putelli/feynman-backend/internal/auth"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"go.uber.org/zap"
)

func handler(handleFunc auth.AuthRefreshTokenHandler) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		refreshToken := request.Headers["Authorization"]
		if refreshToken == "" {
			logger.Error("Missing refresh token")
			return *failure.NewBadRequest("Missing refresh token"), nil
		}

		context, err := auth.NewBaseContext()
		if err != nil {
			logger.Error("Failed to create context", zap.Error(err))
			return *failure.NewInternalServerError(), nil
		}

		tokenInfo, err := handleFunc(context, refreshToken)
		if err != nil {
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
	lambda.Start(handler(auth.HandleRefreshToken))
}
