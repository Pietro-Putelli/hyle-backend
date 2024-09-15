package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/sns"
	"github.com/pietro-putelli/feynman-backend/internal/user"
	"go.uber.org/zap"
)

func handler(ctx context.Context, event events.EventBridgeEvent) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	userContext, err := user.NewContext()
	if err != nil {
		logger.Error("Error creating user context", zap.Error(err))
		return err
	}

	database := userContext.Database

	sessions := []domain.ShortSession{}
	err = database.Model(&domain.Session{}).
		Select("sessions.guid, sessions.user_id, sessions.device_token, users.settings").
		Joins("JOIN users ON users.id = sessions.user_id").
		Where("sessions.expired_at = ? AND (device_token = '') IS NOT TRUE", "0001-01-01 00:00:00").
		Where("users.is_notification_enabled = ?", true).
		Find(&sessions).Error

	if err != nil {
		logger.Error("Error fetching eligible users", zap.Error(err))
		return err
	}

	for _, session := range sessions {
		err = sns.SendMessage(sns.TopicNames.PushNotification, session)
		if err != nil {
			logger.Error("Error sending SNS message", zap.Error(err))
		}
	}

	if err != nil {
		logger.Error("Error fetching eligible users", zap.Error(err))
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
