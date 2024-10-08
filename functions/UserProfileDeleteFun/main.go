package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"github.com/pietro-putelli/feynman-backend/internal/user"
	"github.com/pietro-putelli/feynman-backend/internal/utility"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	userID := utility.GetUserIDBy(request)

	ctx, err := user.NewContext()
	if err != nil {
		return *failure.NewInternalServerError(), nil
	}

	if err := ctx.Service.DeleteUserProfile(userID); err != nil {
		return *failure.NewInternalServerError(), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil
}

func main() {
	lambda.Start(handler)
}
