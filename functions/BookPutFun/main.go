package main

import (
	"fmt"

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

	body := &domain.EditBookBody{}
	if err := utility.ParseRequestBody(request, body); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		return *failure.NewBadRequest(fmt.Sprintf("Invalid parse request body: %s", err.Error())), nil
	}

	ctx, err := book.NewContext()
	if err != nil {
		logger.Error("Failed to create book context", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	if err := ctx.Service.EditBook(userID, body); err != nil {
		logger.Error("Failed to edit book", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil
}

func main() {
	lambda.Start(handler)
}
