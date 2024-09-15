package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"github.com/pietro-putelli/feynman-backend/internal/user"
	"github.com/pietro-putelli/feynman-backend/internal/utility"
	"go.uber.org/zap"
)

/*
	Invoke this function everytime the user opens the app to verify if the user profile is still valid and to check the status of the user subscription.
*/

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	userID := utility.GetUserIDBy(request)

	ctx, err := user.NewContext()
	if err != nil {
		logger.Error("Failed to create user context", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	userInfo, err := ctx.Service.CheckProfileHealth(userID)
	if err != nil {
		logger.Error("Failed to check profile health", zap.Error(err))
		return *failure.NewBadRequest(err.Error()), nil
	}

	response, err := json.Marshal(userInfo)
	if err != nil {
		logger.Error("Failed to marshal response", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(response),
	}, nil
}

func main() {
	lambda.Start(handler)
}
