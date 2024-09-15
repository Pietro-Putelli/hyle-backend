package failure

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// NewInternalServerError creates a new internal server error response.
func NewInternalServerError() *events.APIGatewayProxyResponse {
	err := NewError(500, "An error occurred while processing the request")

	errMessage, _ := json.Marshal(err)
	return &events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       string(errMessage),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}
