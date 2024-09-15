package auth

import (
	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
)

const (
	Google = "google"
	Apple  = "apple"
)

// ThirdPartyProvider represents the third party provider interface.
//
//go:generate mockgen -source=provider.go -destination=./provider_mock.go -package=auth
type ThirdPartyProvider interface {
	ValidateToken(auth *domain.AuthTokenBody) (*domain.ThirdPartyUser, error)
}

type ProviderHandler func(string, *config.Auth) (ThirdPartyProvider, error)

// NewProvider returns the third party provider.
func NewProvider(p string, config *config.Auth) (ThirdPartyProvider, error) {
	switch p {
	case Google:
		return NewGoogleProvider(config.Google.ClientID)
	case Apple:
		return NewAppleProvider(config.Apple)
	default:
		return nil, &InvalidProviderError{}
	}
}

//-------------------------------------
// Errors
//-------------------------------------

// InvalidProviderError represents an error when the provider is not valid.
type InvalidProviderError struct{}

func (e *InvalidProviderError) Error() string {
	return "provider is not valid"
}
