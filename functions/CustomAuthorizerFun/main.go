package main

import (
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pietro-putelli/feynman-backend/internal/auth"
	"go.uber.org/zap"
)

// https://repost.aws/knowledge-center/api-gateway-lambda-authorization-errors

// Help function to generate an IAM policy
func generatePolicy(principalId, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: "user"}

	authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   effect,
				Resource: []string{resource},
			},
		},
	}

	authResponse.Context = map[string]interface{}{
		"userID": principalId,
	}

	return authResponse
}

func handleRequest(event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	raw_token := event.AuthorizationToken
	token := strings.TrimPrefix(raw_token, "Bearer ")

	tmp := strings.Split(event.MethodArn, ":")
	apiGatewayArnTmp := strings.Split(tmp[5], "/")
	resource := tmp[0] + ":" + tmp[1] + ":" + tmp[2] + ":" + tmp[3] + ":" + tmp[4] + ":" + apiGatewayArnTmp[0] + "/*/*"

	context, err := auth.NewBaseContext()

	if err != nil {
		logger.Error("NewBaseContext", zap.Error(err))
		return generatePolicy("User", "Deny", resource), nil
	}

	deniedPolicy := generatePolicy("User", "Deny", resource)

	claims := jwt.MapClaims{}

	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(context.Config.Auth.Jwt.Secret), nil
	})
	if err != nil {
		logger.Error("ParseWithClaims", zap.Error(err))
		return deniedPolicy, nil
	}

	// Check if the token is not of refresh kind
	if claims["kind"] == "refresh" {
		logger.Error("Invalid Token Kind", zap.Error(err))
		return deniedPolicy, nil
	}

	userGuid := claims["sub"].(string)

	return generatePolicy(userGuid, "Allow", resource), nil
}

func main() {
	lambda.Start(handleRequest)
}
