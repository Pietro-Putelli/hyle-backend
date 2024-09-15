package auth

import (
	"context"
	"fmt"

	"github.com/Timothylock/go-signin-with-apple/apple"
	"github.com/mitchellh/mapstructure"
	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
)

var _ ThirdPartyProvider = (*AppleProvider)(nil)

// AppleProvider represents the Apple authentication provider.
type AppleProvider struct {
	clientID     string
	clientSecret string
}

// NewAppleProvider creates a new Apple provider.
func NewAppleProvider(config config.Apple) (*AppleProvider, error) {
	// Generate the client secret used to authenticate with Apple's validation servers
	secret, err := apple.GenerateClientSecret(config.SignInCertificate, config.TeamId, config.AppBundleId, config.SignInCertificateKey)

	if err != nil {
		return nil, err
	}

	return &AppleProvider{
		clientID:     config.AppBundleId,
		clientSecret: secret,
	}, nil
}

// ValidateToken validates the token.
func (a *AppleProvider) ValidateToken(auth *domain.AuthTokenBody) (*domain.ThirdPartyUser, error) {
	client := apple.New()

	vReq := apple.AppValidationTokenRequest{
		ClientID:     a.clientID,
		ClientSecret: a.clientSecret,
		Code:         auth.Token,
	}

	var resp apple.ValidationResponse

	// Do the verification
	err := client.VerifyAppToken(context.Background(), vReq, &resp)

	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		fmt.Printf("[Apple Provider Error]: %s - %s\n", resp.Error, resp.ErrorDescription)
		return nil, err
	}

	// Get the unique user ID
	unique, err := apple.GetUniqueID(resp.IDToken)
	if err != nil {
		return nil, err
	}

	// Get the email
	claim, err := apple.GetClaims(resp.IDToken)
	if err != nil {
		return nil, err
	}

	email := (*claim)["email"]

	if email == nil {
		return nil, fmt.Errorf("[Apple Provider Error]: Email not found in claims")
	}

	thirdPartyUser := domain.NewAppleThirdPartyUser()

	if err := mapstructure.Decode(
		map[string]interface{}{
			"email":        email,
			"given_name":   auth.Data.GivenName,
			"family_name":  auth.Data.FamilyName,
			"sub":          unique,
			"device_id":    auth.Device.ID,
			"device_token": auth.Device.Token,
		}, thirdPartyUser); err != nil {
		return nil, err
	}

	return thirdPartyUser, nil
}
