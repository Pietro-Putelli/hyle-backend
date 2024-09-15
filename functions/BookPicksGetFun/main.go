package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/pietro-putelli/feynman-backend/internal/book"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"github.com/pietro-putelli/feynman-backend/internal/utility"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	params := &domain.GetPicksParams{}
	utility.ParseQueryParams(request, params)

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(params); err != nil {
		logger.Error("Params validation failed", zap.Error(err))
		return *failure.NewBadRequest("Params validation failed"), nil
	}

	ctx, err := book.NewContext()
	if err != nil {
		logger.Error("Failed to create context", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	picks, err := ctx.Service.GetPicksByBook(params)
	if err != nil {
		logger.Error("Failed to get picks by book", zap.Error(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return *failure.NewNotFound("Book not found"), nil
		}

		return *failure.NewInternalServerError(), nil
	}

	response, err := json.Marshal(picks)
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
