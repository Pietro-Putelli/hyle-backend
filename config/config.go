package config

import (
	"errors"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config represents the configuration of the application.
	Config struct {
		// Auth represents the authentication configuration.
		Auth Auth
		// Database represents the database configuration.
		Database Database
		// Apple represents the Apple configuration.
		Apple Apple
		// Langchain represents the langchain configuration.
		Langchain Langchain

		// Telegram represents the Telegram configuration.
		Telegram Telegram
	}

	// Auth represents the authentication configuration.
	Auth struct {
		Google AuthGoogle
		Apple  Apple
		Jwt    AuthJwt
	}

	// AuthGoogle represents the Google authentication configuration.
	AuthGoogle struct {
		ClientID string `env-required:"true" env:"AUTH_GOOGLE_CLIENT_ID"`
	}

	// Apple represents the Apple configuration.
	Apple struct {
		TeamId      string `env-required:"true" env:"APPLE_TEAM_ID"`
		AppBundleId string `env-required:"true" env:"IOS_APP_BUNDLE_ID"`
		IssuerId    string `env-required:"true" env:"APPSTORE_ISSUER_ID"`

		SignInCertificate    string `env-required:"true" env:"APPLE_SIGNIN_CERTIFICATE"`
		SignInCertificateKey string `env-required:"true" env:"APPLE_SIGNIN_CERTIFICATE_KEY"`

		ApnsCertificate    string `env-required:"true" env:"APPLE_APNS_CERTIFICATE"`
		ApnsCertificateKey string `env-required:"true" env:"APPLE_APNS_CERTIFICATE_KEY"`

		// IAPCertificate    string `env-required:"true" env:"APPLE_IAP_CERTIFICATE"`
		// IAPCertificateKey string `env-required:"true" env:"APPLE_IAP_CERTIFICATE_KEY"`
	}

	// AuthJwt represents the JWT authentication configuration.
	AuthJwt struct {
		Secret              string `env-required:"true" env:"JWT_SECRET"`
		AccessTokenDuration int    `env-default:"604800" env:"JWT_ACCESS_TOKEN_DURATION"`
	}

	// Database represents the database configuration.
	Database struct {
		Host     string `env-required:"true" env:"DB_HOST"`
		Port     int    `env-required:"true" env:"DB_PORT"`
		User     string `env-required:"true" env:"DB_USER"`
		Password string `env-required:"true" env:"DB_PASSWORD"`
		Name     string `env-required:"true" env:"DB_NAME"`
		SSLMode  string `env-default:"disable" env:"DB_SSL_MODE"`
	}

	// Langchain represents the langchain configuration.
	Langchain struct {
		OpenAIKey    string `env-required:"true" env:"OPENAI_API_KEY"`
		LangsmithKey string `env-required:"true" env:"LANGCHAIN_API_KEY"`
		GPTModel     string `env-required:"true" env:"GPT_MODEL"`
	}

	Telegram struct {
		ApiToken string `env-required:"true" env:"TELEGRAM_API_TOKEN"`
	}
)

// NewConfig creates a new configuration.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, errors.New("failed to load config: " + err.Error())
	}

	return cfg, nil
}
