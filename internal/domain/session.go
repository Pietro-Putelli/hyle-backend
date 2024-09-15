package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

//-------------------------------------
// Session DB Model
//-------------------------------------

// Session represents the session domain.
type Session struct {
	TimestapModel

	Guid        uuid.UUID `gorm:"type:uuid;unique;not null;column:guid;default:uuid_generate_v4()"`
	User        User      `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE"`
	UserID      uint      `gorm:"column:user_id;not null"`
	DeviceID    string    `gorm:"column:device_id;not null"`
	DeviceToken string    `gorm:"column:device_token"`
	ExpiredAt   time.Time `gorm:"column:expired_at"`
}

type ShortSession struct {
	Guid        uuid.UUID       `json:"guid"`
	UserID      uint            `json:"user_id"`
	DeviceToken string          `json:"device_token"`
	Settings    json.RawMessage `json:"settings"`
}

func (Session) TableName() string {
	return "sessions"
}

// Session Update Request

type PatchSessionBody struct {
	Guid        string `json:"guid" validate:"required"`
	DeviceToken string `json:"device_token" validate:"required"`
}
