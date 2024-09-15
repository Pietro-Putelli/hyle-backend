package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/failure"
)

type AuthTokenHandler func(*Context, *domain.AuthTokenBody) (*domain.AuthUserTokenDto, error)

// HandleAuthToken handles the authentication token request.
func HandleAuthToken(ctx *Context, body *domain.AuthTokenBody) (*domain.AuthUserTokenDto, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Validate the request body
	if err := validate.Struct(body); err != nil {
		return nil, failure.NewValidationErr(err)
	}

	// Get authentication provider
	provider, err := ctx.ProviderHandler(body.Provider, &ctx.Config.Auth)
	if err != nil {
		return nil, err
	}

	// Validating the token
	thirdPartyUser, err := provider.ValidateToken(body)
	if err != nil {
		return nil, err
	}

	// Create user if it does not exist
	user, sessionID, err := ctx.Service.CreateUserIfNotExists(thirdPartyUser)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	authToken, err := ctx.Service.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthUserTokenDto{
		User: domain.UserDtoFromModel(*user, sessionID),
		Auth: authToken,
	}, nil
}

type AuthRefreshTokenHandler func(*Context, string) (*domain.AuthTokenDto, error)

// HandleRefreshToken handles the refresh token request.
func HandleRefreshToken(ctx *Context, refreshToken string) (*domain.AuthTokenDto, error) {
	// Validate the refresh token
	authToken, err := ctx.Service.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return authToken, nil
}
