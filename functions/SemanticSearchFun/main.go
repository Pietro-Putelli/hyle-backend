package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

	params := &domain.SearchGetParams{}
	utility.ParseQueryParams(request, params)

	userID := utility.GetUserIDBy(request)

	ctx, err := book.NewContext()
	if err != nil {
		logger.Error("Error creating context", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	var bookPicks any

	if params.BookID != "" {
		bookPicks, err = ctx.Service.SearchPickInBook(params)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return *failure.NewNotFound("No results found"), nil
			}

			logger.Error("Error searching picks in book", zap.Error(err))
			return *failure.NewInternalServerError(), nil
		}
	} else {
		bookPicks, err = ctx.Service.SemanticSearch(userID, params)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return *failure.NewNotFound("No results found"), nil
			}

			logger.Error("Error searching semantic", zap.Error(err))
			return *failure.NewInternalServerError(), nil
		}
	}

	response, err := json.Marshal(bookPicks)
	if err != nil {
		logger.Error("Failed to marshal response", zap.Error(err))
		return *failure.NewInternalServerError(), nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(response),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
