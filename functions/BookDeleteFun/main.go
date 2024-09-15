package main

import (
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
		logger.Error("Invalid book guid", zap.Error(err))
		return *failure.NewBadRequest("Invalid Book ID"), nil
	}

	userID := utility.GetUserIDBy(request)

	ctx, err := book.NewContext()
	if err != nil {
		logger.Error("Error creating book context", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	if err := ctx.Service.DeleteBook(userID, bookID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return *failure.NewNotFound("Book not found"), nil
		}

		logger.Error("Invalid book guid", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil
}

func main() {
	lambda.Start(handler)
}
