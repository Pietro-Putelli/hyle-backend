package failure

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// NewBadRequest creates a new bad request response.
func NewBadRequest(message string) *events.APIGatewayProxyResponse {
	err := NewError(400, message)

	errMessage, _ := json.Marshal(err)
	return &events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       string(errMessage),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

// NewNotFound creates a new not found response.
func NewNotFound(message string) *events.APIGatewayProxyResponse {
	err := NewError(404, message)

	errMessage, _ := json.Marshal(err)
	return &events.APIGatewayProxyResponse{
		StatusCode: 404,
		Body:       string(errMessage),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}
