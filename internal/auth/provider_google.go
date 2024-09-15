package auth

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"golang.org/x/net/context"
	"google.golang.org/api/idtoken"
)

var _ ThirdPartyProvider = (*GoogleProvider)(nil)

// GoogleInternalValidator represents the Google internal validator.
type GoogleInternalValidator interface {
	Validate(ctx context.Context, token, audience string) (*idtoken.Payload, error)
}

// GoogleProvider represents the Google authentication provider.
type GoogleProvider struct {
	validator GoogleInternalValidator
	clientID  string
}

// NewGoogleProvider creates a new Google provider.
func NewGoogleProvider(clientID string) (*GoogleProvider, error) {
	validator, err := idtoken.NewValidator(context.Background())
	if err != nil {
		return nil, err
	}
	return &GoogleProvider{
		clientID:  clientID,
		validator: validator,
	}, nil
}

// ValidateToken validates the token.
func (g *GoogleProvider) ValidateToken(auth *domain.AuthTokenBody) (*domain.ThirdPartyUser, error) {
	payload, err := g.validator.Validate(context.Background(), auth.Token, g.clientID)
	if err != nil {
		return nil, err
	}

	// build thirdpartyuser from payload claims
	thirdPartyUser := domain.NewGoogleThirdPartyUser()
	claims := payload.Claims

	if err := mapstructure.Decode(map[string]interface{}{
		"email":        claims["email"],
		"given_name":   claims["given_name"],
		"family_name":  claims["family_name"],
		"sub":          claims["sub"],
		"device_id":    auth.Device.ID,
		"device_token": auth.Device.Token,
	}, thirdPartyUser); err != nil {
		return nil, err
	}
	return thirdPartyUser, nil
}
