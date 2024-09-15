package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/pietro-putelli/feynman-backend/internal/book"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"github.com/pietro-putelli/feynman-backend/internal/utility"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	stringBookID := request.PathParameters["bookId"]

	bookID, err := uuid.Parse(stringBookID)
	if err != nil {
		logger.Error("Invalid Book ID", zap.Error(err))
		return *failure.NewBadRequest("Invalid Book ID"), nil
	}

	userID := utility.GetUserIDBy(request)

	ctx, err := book.NewContext()
	if err != nil {
		logger.Error("Failed to create context", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	bookResponse, err := ctx.Service.GetCompleteBookByGuid(userID, bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return *failure.NewNotFound("Book not found"), nil
		}

		logger.Error("Failed to get book", zap.Error(err))

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return *failure.NewNotFound("Book not found"), nil
		}

		return *failure.NewInternalServerError(), nil
	}

	response, err := json.Marshal(bookResponse)
	if err != nil {
		logger.Error("Failed to marshal single book response", zap.Error(err))
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

func main() {
	lambda.Start(handler)
}
