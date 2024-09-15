package domain

import (
	"time"
)

//-------------------------------------
// AuthToken Response
//-------------------------------------

// AuthTokenDto represents the authentication token domain.
type AuthTokenDto struct {
	AccessToken          string    `json:"access_token"`
	RefreshToken         string    `json:"refresh_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

// AuthUserTokenDto represents the authentication user with token data transfer object.
type AuthUserTokenDto struct {
	User *UserDto      `json:"user"`
	Auth *AuthTokenDto `json:"auth"`
}

// AuthTokenDevice used when user is created. For patch see /session/PatchSessionBody
type AuthTokenDevice struct {
	// Unique identifier of the mobile device
	ID string `json:"id"`
	// Mobile Token to be used for push notifications
	Token string `json:"token"`
}

//-------------------------------------
// AuthToken Request
//-------------------------------------

// AuthTokenBody represents the authentication token body.
type AuthTokenBody struct {
	Token    string           `json:"token" validate:"required"`
	Provider string           `json:"provider" validate:"required,oneof=apple google" enum:"apple,google"`
	Data     *AuthTokenData   `json:"data" validate:"required_if=Provider apple"`
	Device   *AuthTokenDevice `json:"device" validate:"required"`
}

// AuthTokenData represents the authentication token data.
type AuthTokenData struct {
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}
