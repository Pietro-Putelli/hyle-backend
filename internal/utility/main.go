package utility

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	IsLocalEnv bool `env:"IS_LOCAL_ENV" env-default:"false"`
}

// GetUserIDBy returns the userID from the request
func GetUserIDBy(request events.APIGatewayProxyRequest) uuid.UUID {
	config := &Config{}
	err := cleanenv.ReadEnv(config)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err != nil {
		logger.Error("Failed to read config", zap.Error(err))
	}

	stringUserID := ""

	if config.IsLocalEnv {
		stringUserID = "16bebb13-2dfa-4137-918d-be3aa3ef940a"
	} else {
		stringUserID = request.RequestContext.Authorizer["userID"].(string)
	}

	userID, err := uuid.Parse(stringUserID)

	if err != nil {
		logger.Error("Failed to parse userID", zap.String("userID", stringUserID), zap.Error(err))
		return uuid.UUID{}
	}

	return userID
}
