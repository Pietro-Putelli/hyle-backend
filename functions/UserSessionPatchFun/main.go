package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"github.com/pietro-putelli/feynman-backend/internal/session"
	"github.com/pietro-putelli/feynman-backend/internal/utility"
	"go.uber.org/zap"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	userID := utility.GetUserIDBy(request)

	body := &domain.PatchSessionBody{}
	if err := json.Unmarshal([]byte(request.Body), body); err != nil {
		return *failure.NewBadRequest("Invalid request body"), nil
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(body); err != nil {
		return *failure.NewBadRequest("Invalid request body"), nil
	}

	ctx, err := session.NewContext()
	if err != nil {
		return *failure.NewInternalServerError(), nil
	}

	if err := ctx.Service.UpdateUserSession(userID, body); err != nil {
		logger.Error("Failed to update user session", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil
}

func main() {
	lambda.Start(handler)
}
