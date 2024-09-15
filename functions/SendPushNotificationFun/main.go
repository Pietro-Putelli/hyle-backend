package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/user"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
	"go.uber.org/zap"
)

type SnsMessage struct {
	Message string `json:"message"`
}

func handler(ctx context.Context, event events.SNSEvent) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	message := event.Records[0].SNS.Message

	userSession := domain.ShortSession{}

	err := json.Unmarshal([]byte(message), &userSession)
	if err != nil {
		logger.Error("Error unmarshalling user", zap.Error(err))
		return err
	}

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("Error loading config", zap.Error(err))
		return err
	}

	userContext, err := user.NewContext()
	if err != nil {
		logger.Error("Error creating user context", zap.Error(err))
		return err
	}

	database := userContext.Database

	userSettings := domain.UserSettings{}
	err = json.Unmarshal(userSession.Settings, &userSettings)
	if err != nil {
		logger.Error("Error unmarshalling user settings", zap.Error(err))
		return err
	}

	limit := 3
	if userSettings.NotificationMode == "all" {
		limit = -1
	}

	var pick domain.BookPickPushNotification

	err = database.Model(&domain.Book{}).
		Select("books.guid AS book_id, book_picks.content_text AS content").
		Where("books.user_id = ?", userSession.UserID).
		Order("books.updated_at DESC").
		Limit(limit).
		Joins("JOIN book_picks ON books.id = book_picks.book_id").
		Order("RANDOM()").
		First(&pick).Error

	if err != nil {
		logger.Error("Error fetching book pick", zap.Error(err))
		return err
	}

	apnsKey, err := token.AuthKeyFromBytes([]byte(cfg.Apple.ApnsCertificate))
	if err != nil {
		logger.Error("Error loading APNS key", zap.Error(err))
		return err
	}

	token := &token.Token{
		AuthKey: apnsKey,
		KeyID:   cfg.Apple.ApnsCertificateKey,
		TeamID:  cfg.Apple.TeamId,
	}

	payload := domain.PushNotificationPayload{
		Aps: domain.PusNotificationAps{
			Alert: domain.PushNotificationAlert{
				Title: "📚 Daily Pick Reminder",
				Body:  pick.Content,
			},
			Badge: 1,
		},
		Data: domain.PushNotificationPayloadData{
			BookID: pick.BookID,
		},
	}

	notification := &apns2.Notification{
		DeviceToken: userSession.DeviceToken,
		Topic:       cfg.Apple.AppBundleId,
		Payload:     payload,
	}

	client := apns2.NewTokenClient(token)
	client.Host = apns2.HostProduction

	_, err = client.Push(notification)
	if err != nil {
		logger.Error("Error sending push notification", zap.Error(err))
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
