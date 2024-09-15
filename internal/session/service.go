package session

import (
	"time"

	"github.com/google/uuid"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/user"
	"gorm.io/gorm"
)

var _ Service = (*serviceImpl)(nil)

// The Service interface is defined to abstract the behaviors that the serviceImpl struct will implement.
type Service interface {
	UpdateUserSession(userID uuid.UUID, session *domain.PatchSessionBody) error
	LogoutSession(userID uuid.UUID, sessionID uuid.UUID) error
}

type serviceImpl struct {
	db          *gorm.DB
	userService user.Service
}

func NewService(db *gorm.DB, userService user.Service) Service {
	return &serviceImpl{
		db:          db,
		userService: userService,
	}
}

func (service *serviceImpl) UpdateUserSession(userID uuid.UUID, s *domain.PatchSessionBody) error {
	user, err := service.userService.GetUserByGuid(userID)
	if err != nil {
		return err
	}

	sessionID := s.Guid

	data := map[string]interface{}{
		"updated_at": time.Now(),
	}

	return service.db.Transaction(func(tx *gorm.DB) error {

		if s.DeviceToken != "" {
			data["device_token"] = s.DeviceToken

			/* If the device token is updated, the user has granted permission for push notifications, therefore update corresponding user entry */
			err := tx.Model(&domain.User{}).Where("id = ?", user.ID).Update("is_notification_enabled", true).Error
			if err != nil {
				return err
			}
		}

		// Search for the session in the database and update it
		if err := tx.Model(&domain.Session{}).Where("user_id = ? AND guid = ?", user.ID, sessionID).Updates(data).Error; err != nil {
			return err
		}

		return nil
	})
}

func (service *serviceImpl) LogoutSession(userID uuid.UUID, sessionID uuid.UUID) error {
	user, err := service.userService.GetUserByGuid(userID)

	if err != nil {
		return err
	}

	// Search for the session in the database and delete it
	if err := service.db.Model(&domain.Session{}).Where("user_id = ? AND guid = ?", user.ID, sessionID).Update("expired_at", time.Now()).Error; err != nil {
		return err
	}

	return nil
}
