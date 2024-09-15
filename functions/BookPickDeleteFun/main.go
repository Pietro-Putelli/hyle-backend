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

	params := &domain.DeleteBookPickPath{}
	if err := utility.ParsePathParams(request, params); err != nil {
		logger.Error("Error parsing path params", zap.Error(err))
		return *failure.NewBadRequest("Error parsing path params"), nil
	}

	userID := utility.GetUserIDBy(request)

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(params); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		return *failure.NewBadRequest("Params validation failed"), nil
	}

	ctx, err := book.NewContext()
	if err != nil {
		logger.Error("Error creating book context", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	isLastPick, err := ctx.Service.DeleteBookPick(userID, params)
	if err != nil {
		logger.Error("Error deleting book pick", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	response := map[string]interface{}{
		"is_last": isLastPick,
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		logger.Error("Error serializing response", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	logger.Info("Book pick deleted", zap.Any("response", response))
	logger.Info("Book pick deleted", zap.Any("response", string(responseBody)))

	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
