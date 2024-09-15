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

	body := &domain.EditBookPickBody{}
	if err := json.Unmarshal([]byte(request.Body), body); err != nil {
		logger.Error("Invalid Request Body", zap.Error(err))
		return *failure.NewBadRequest("Invalid Request Body"), nil
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(body); err != nil {
		logger.Error("Invalid Request Body", zap.Error(err))
		return *failure.NewBadRequest("Body Validation Failed"), nil
	}

	ctx, err := book.NewContext()
	if err != nil {
		logger.Error("Failed to Create Context", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	if err = ctx.Service.EditBookPick(userID, body); err != nil {
		logger.Error("Unable to Update Pick", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil
}

func main() {
	lambda.Start(handler)
}
