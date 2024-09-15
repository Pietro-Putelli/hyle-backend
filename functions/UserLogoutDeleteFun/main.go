package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"github.com/pietro-putelli/feynman-backend/internal/session"
	"github.com/pietro-putelli/feynman-backend/internal/utility"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	userID := utility.GetUserIDBy(request)

	rawSessionID := request.QueryStringParameters["guid"]

	if rawSessionID == "" {
		return *failure.NewBadRequest("Invalid Session ID"), nil
	}

	ctx, err := session.NewContext()
	if err != nil {
		return *failure.NewInternalServerError(), nil
	}

	sessionID, err := uuid.Parse(rawSessionID)
	if err != nil {
		return *failure.NewBadRequest("Invalid Session ID"), nil
	}

	if err := ctx.Service.LogoutSession(userID, sessionID); err != nil {
		return *failure.NewInternalServerError(), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil
}

func main() {
	lambda.Start(handler)
}
