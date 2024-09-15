package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
	"github.com/pietro-putelli/feynman-backend/internal/user"
	"github.com/pietro-putelli/feynman-backend/internal/utility"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	userID := utility.GetUserIDBy(request)

	body := &domain.UserProfileUpdate{}
	if err := json.Unmarshal([]byte(request.Body), body); err != nil {
		return *failure.NewBadRequest("Invalid request body"), nil
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(body); err != nil {
		return *failure.NewBadRequest("Invalid request body"), nil
	}

	ctx, err := user.NewContext()
	if err != nil {
		return *failure.NewInternalServerError(), nil
	}

	if err := ctx.Service.UpdateUserProfile(userID, body); err != nil {
		return *failure.NewBadRequest(err.Error()), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil
}

func main() {
	lambda.Start(handler)
}
