package user

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/utility"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ Service = (*serviceImpl)(nil)

// Service represents the user service.
//
//go:generate mockgen -source=service.go -destination=./service_mock.go -package=user
type Service interface {
	GetUserByGuid(guid uuid.UUID) (*domain.User, error)
	CreateUserIfNotExists(user *domain.ThirdPartyUser) (*domain.User, uuid.UUID, error)
	UpdateUserProfile(userID uuid.UUID, data *domain.UserProfileUpdate) error
	CheckProfileHealth(userID uuid.UUID) (*domain.UserHealth, error)
}

// serviceImpl represents the user service implementation.
type serviceImpl struct {
	db *gorm.DB
}

// NewService creates a new user service.
func NewService(db *gorm.DB) Service {
	return &serviceImpl{
		db: db,
	}
}

// GetUserByGuid retrieves a user by its GUID.
func (s *serviceImpl) GetUserByGuid(guid uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := s.db.Where("guid = ?", guid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUserIfNotExists creates a new user if it does not exist.
func (s *serviceImpl) CreateUserIfNotExists(user *domain.ThirdPartyUser) (*domain.User, uuid.UUID, error) {
	var responseUser domain.User
	var sessionID = uuid.UUID{}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.db.Where("email = ?", user.Email).First(&responseUser).Error; err != nil {
			responseUser = domain.User{
				Email:      user.Email,
				GivenName:  user.GivenName,
				FamilyName: user.FamilyName,
				ExternalID: user.Sub,
				Provider:   user.Provider,
				Settings:   domain.NewUserSettings(),
				IsCreated:  true,
			}

			if err := s.db.Create(&responseUser).Error; err != nil {
				logger.Error("Failed to create user", zap.Error(err))
				return err
			}

			/* New user has been created, copy the sample books to gettings started */
			err := utility.CopyBookSamplesToUser(responseUser.ID, tx)
			if err != nil {
				logger.Error("Failed to copy book samples", zap.Error(err))
				return err
			}

			utility.TelegramSendNewUser(&responseUser)

			newSession := domain.Session{
				UserID:   responseUser.ID,
				DeviceID: user.DeviceID,
				// DeviceToken can be empty (not null) if the user has not granted notification permission
				DeviceToken: user.DeviceToken,
			}

			if err := s.db.Create(&newSession).Error; err != nil {
				logger.Error("Failed to create new session", zap.Error(err))
				return err
			}

			sessionID = newSession.Guid

			return nil
		}

		/* If the user already exists, update existing and not expired session or create a new session */
		var existingSession domain.Session
		err := s.db.Model(&domain.Session{}).Where("user_id = ? AND expired_at = ?", responseUser.ID, "0001-01-01 00:00:00").First(&existingSession).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newSession := domain.Session{
					UserID:      responseUser.ID,
					DeviceID:    user.DeviceID,
					DeviceToken: user.DeviceToken,
				}

				if err := s.db.Create(&newSession).Error; err != nil {
					logger.Error("Failed to create new session", zap.Error(err))
					return err
				}

				sessionID = newSession.Guid

			} else {
				logger.Error("Failed to get existing session", zap.Error(err))
				return err
			}
		} else {
			updateData := map[string]interface{}{
				"device_id":    user.DeviceID,
				"device_token": user.DeviceToken,
			}

			if err := s.db.Model(&domain.Session{}).Where("guid = ?", existingSession.Guid).Updates(updateData).Error; err != nil {
				logger.Error("Failed to update existing session", zap.Error(err))
				return err
			}

			sessionID = existingSession.Guid
		}

		return nil
	})

	return &responseUser, sessionID, err
}

// UpdateUserProfile changes a restricted set of user profile fields.

func (s *serviceImpl) UpdateUserProfile(userID uuid.UUID, data *domain.UserProfileUpdate) error {

	fieldsToUpdate := map[string]interface{}{}

	jsonSettings, err := json.Marshal(data.Settings)
	if err == nil {
		fieldsToUpdate["settings"] = string(jsonSettings)
		fieldsToUpdate["is_notification_enabled"] = data.Settings.NotificationEnabled
	}

	if data.SubscriptionReceiptID != "" {
		fieldsToUpdate["subscription_receipt_id"] = data.SubscriptionReceiptID
	}

	if err := s.db.Model(&domain.User{}).Where("guid = ?", userID).Updates(fieldsToUpdate); err != nil {
		return err.Error
	}

	return nil
}

// CheckProfileHealth checks the user is still valid and the user subscription status.
func (s *serviceImpl) CheckProfileHealth(userID uuid.UUID) (*domain.UserHealth, error) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	user, err := s.GetUserByGuid(userID)
	if err != nil {
		return nil, err
	}

	response := &domain.UserHealth{
		IsPremium: true,
		IsHealthy: false,
	}

	if !user.IsActive {
		return response, nil
	}

	response.IsHealthy = true

	// if user.SubscriptionReceiptID != "" {

	// 	/* Check the receipt validity against Apple Server */
	// 	ctx, err := NewContext()
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	appleConfig := ctx.Config.Apple

	// 	// https://github.com/awa/go-iap?tab=readme-ov-file#in-app-store-server-api

	// 	config := &api.StoreConfig{
	// 		KeyContent: []byte(appleConfig.IAPCertificate),
	// 		KeyID:      appleConfig.IAPCertificateKey,
	// 		BundleID:   appleConfig.AppBundleId,
	// 		Issuer:     appleConfig.IssuerId,
	// 		Sandbox:    true,
	// 	}

	// 	transactionId := user.SubscriptionReceiptID
	// 	storeClient := api.NewStoreClient(config)

	// 	ctx1 := context.Background()

	// 	transactionResponse, err := storeClient.GetTransactionInfo(ctx1, transactionId)
	// 	if err != nil {
	// 		logger.Error("Failed to get transaction info", zap.Error(err))
	// 		return nil, err
	// 	}

	// 	transaction, err := storeClient.ParseSignedTransaction(transactionResponse.SignedTransactionInfo)
	// 	if err != nil {
	// 		logger.Error("Failed to parse signed transaction", zap.Error(err))
	// 		return nil, err
	// 	}

	// 	if transaction.OriginalTransactionId == transactionId {
	// 		response.IsPremium = true
	// 	}
	// }

	return response, nil
}
