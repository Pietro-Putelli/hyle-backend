package domain

import (
	"github.com/google/uuid"
)

//-------------------------------------
// User DB Model
//-------------------------------------

// User represents the user domain.
type User struct {
	TimestapModel

	ID         uint          `gorm:"primaryKey;autoIncrement;column:id"`
	Guid       uuid.UUID     `gorm:"type:uuid;unique;not null;column:guid;default:uuid_generate_v4()"`
	ExternalID string        `gorm:"unique;not null;column:external_id"`
	Email      string        `gorm:"unique;not null;column:email"`
	GivenName  string        `gorm:"column:given_name;not null"`
	FamilyName string        `gorm:"column:family_name;not null"`
	Provider   string        `gorm:"column:provider;not null"`
	Settings   *UserSettings `gorm:"column:settings;type:jsonb;serializer:json"`

	SubscriptionReceiptID string `gorm:"column:subscription_receipt_id"`

	IsNotificationEnabled bool `gorm:"column:is_notification_enabled;default:false"`
	IsActive              bool `gorm:"column:is_active;default:true"`

	IsCreated bool `gorm:"-:all"`
}

// TableName returns the table name for the user domain.
func (User) TableName() string {
	return "users"
}

// UserSettings represents the user settings domain.
type UserSettings struct {
	DarkMode       bool   `json:"darkMode"`
	AppLanguage    string `json:"appLanguage"`
	SecondLanguage string `json:"secondLanguage"`

	NotificationEnabled bool `json:"notificationEnabled"`
	/* all, last-edit */
	NotificationMode string `json:"notificationMode"`
}

// NewUserSettings creates a new user settings.
func NewUserSettings() *UserSettings {
	return &UserSettings{
		DarkMode:            true,
		AppLanguage:         "",
		SecondLanguage:      "",
		NotificationEnabled: true,
		NotificationMode:    "all",
	}
}

//-------------------------------------
// User DTO
//-------------------------------------

// UserDto represents the user data transfer object.
type UserDto struct {
	Guid       uuid.UUID     `json:"guid"`
	Email      string        `json:"email"`
	Settings   *UserSettings `json:"settings"`
	GivenName  string        `json:"given_name"`
	FamilyName string        `json:"family_name"`
	IsCreated  bool          `json:"is_created"`
	SessionID  uuid.UUID     `json:"session_id"`
	IsPremium  bool          `json:"is_premium"`
}

// UserDtoFromModel converts a user model to a user data transfer object.
func UserDtoFromModel(user User, sessionID uuid.UUID) *UserDto {
	return &UserDto{
		Guid:       user.Guid,
		Email:      user.Email,
		GivenName:  user.GivenName,
		FamilyName: user.FamilyName,
		Settings:   user.Settings,
		IsCreated:  user.IsCreated,
		IsPremium:  user.SubscriptionReceiptID != "",
		SessionID:  sessionID,
	}
}

//-------------------------------------
// User Third Party
//-------------------------------------

// ThirdPartyUser represents the third party user domain.
type ThirdPartyUser struct {
	Email      string `mapstructure:"email"`
	Sub        string `mapstructure:"sub"` // external id
	Name       string `mapstructure:"name"`
	GivenName  string `mapstructure:"given_name"`
	FamilyName string `mapstructure:"family_name"`
	Provider   string

	DeviceID    string `mapstructure:"device_id"`
	DeviceToken string `mapstructure:"device_token"`
}

// NewGoogleThirdPartyUser creates a new Google third party user.
func NewGoogleThirdPartyUser() *ThirdPartyUser {
	return &ThirdPartyUser{
		Provider: "google",
	}
}

// NewAppleThirdPartyUser creates a new Apple third party user.
func NewAppleThirdPartyUser() *ThirdPartyUser {
	return &ThirdPartyUser{
		Provider: "apple",
	}
}

// UserProfileUpdate represents the user profile update domain.
type UserProfileUpdate struct {
	SubscriptionReceiptID string       `json:"subscription_receipt_id"`
	Settings              UserSettings `json:"settings"`
}

// UserHealth body response
type UserHealth struct {
	IsHealthy bool `json:"is_healthy"`
	IsPremium bool `json:"is_premium"`
}
