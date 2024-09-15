package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/user"
)

// makes sure to implement the Service interface
var _ Service = (*serviceImpl)(nil)

// Service represents the authentication service.
//
//go:generate mockgen -source=service.go -destination=./service_mock.go -package=auth
type Service interface {
	GenerateToken(user *domain.User) (*domain.AuthTokenDto, error)
	RefreshAccessToken(refreshToken string) (*domain.AuthTokenDto, error)
	CreateUserIfNotExists(user *domain.ThirdPartyUser) (*domain.User, uuid.UUID, error)
}

// serviceImpl represents the authentication service implementation.
type serviceImpl struct {
	jwtSecret    []byte
	jwtExpiresIn int
	userService  user.Service
}

// NewService creates a new authentication service.
func NewService(jwtConfig *config.AuthJwt, userService user.Service) Service {
	return &serviceImpl{
		userService:  userService,
		jwtSecret:    []byte(jwtConfig.Secret),
		jwtExpiresIn: jwtConfig.AccessTokenDuration,
	}
}

// generateAccessToken generates a new access token.
func (s *serviceImpl) generateAccessToken(user *domain.User, expirationDate *time.Time) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.Guid,
		"email": user.Email,
		"exp":   expirationDate.Unix(),
		"iat":   time.Now().Unix(),
		"iss":   "panta.srvless-api",
		"kind":  "access",
	})
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

// GenerateToken generates a new token validating the token ID.
func (s *serviceImpl) GenerateToken(user *domain.User) (*domain.AuthTokenDto, error) {
	// Generate access token
	accessTokenExpiresAt := time.Now().Add(time.Duration(s.jwtExpiresIn) * time.Minute)
	accessTokenString, err := s.generateAccessToken(user, &accessTokenExpiresAt)
	if err != nil {
		return nil, err
	}

	// Generate infinite refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.Guid,
		"email": user.Email,
		"exp":   time.Now().AddDate(1000, 0, 0).Unix(),
		"iat":   time.Now().Unix(),
		"iss":   "panta.srvless-api",
		"kind":  "refresh",
	})
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &domain.AuthTokenDto{
		AccessToken:          accessTokenString,
		RefreshToken:         refreshTokenString,
		AccessTokenExpiresAt: accessTokenExpiresAt,
	}, nil
}

// RefreshAccessToken generates a new access token using the refresh token.
func (s *serviceImpl) RefreshAccessToken(refreshToken string) (*domain.AuthTokenDto, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Check if the token is a refresh token
	if claims["kind"] != "refresh" {
		return nil, errors.New("invalid token")
	}

	// Get user by GUID
	userGuidStr, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid token subject")
	}

	// Parse GUID
	userGuid, err := uuid.Parse(userGuidStr)
	if err != nil {
		return nil, errors.New("invalid token subject")
	}

	// Get user
	user, err := s.userService.GetUserByGuid(userGuid)
	if err != nil {
		return nil, err
	}

	// Generate access token
	accessTokenExpiresAt := time.Now().Add(time.Duration(s.jwtExpiresIn) * time.Minute)
	accessTokenString, err := s.generateAccessToken(user, &accessTokenExpiresAt)
	if err != nil {
		return nil, err
	}

	return &domain.AuthTokenDto{
		AccessToken:          accessTokenString,
		RefreshToken:         refreshToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
	}, nil
}

// CreateUserIfNotExists creates a new user if it does not exist.
func (s *serviceImpl) CreateUserIfNotExists(user *domain.ThirdPartyUser) (*domain.User, uuid.UUID, error) {
	return s.userService.CreateUserIfNotExists(user)
}
