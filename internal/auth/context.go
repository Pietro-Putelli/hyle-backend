package auth

import (
	"errors"

	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/pietro-putelli/feynman-backend/internal/database"
	"github.com/pietro-putelli/feynman-backend/internal/user"
	"gorm.io/gorm"
)

// Context represents the context of the authentication service.
type Context struct {
	Service         Service
	UserService     user.Service
	Config          *config.Config
	DB              *gorm.DB
	ProviderHandler ProviderHandler
}

// NewContext creates a new authentication context.
func NewContext(cfg *config.Config, db *gorm.DB, providerHandler ProviderHandler, service Service) *Context {
	return &Context{
		Service:         service,
		UserService:     user.NewService(db),
		Config:          cfg,
		DB:              db,
		ProviderHandler: providerHandler,
	}
}

// NewBaseContext creates a new base context.
func NewBaseContext() (*Context, error) {
	// load configuration
	config, err := config.NewConfig()
	if err != nil {
		return nil, errors.New("failed load auth context config: " + err.Error())
	}

	// load database
	database, err := database.NewDB(database.NewConn(&config.Database))
	if err != nil {
		return nil, errors.New("failed load auth context database: " + err.Error())
	}

	// user service
	userService := user.NewService(database)

	return NewContext(config, database, NewProvider, NewService(&config.Auth.Jwt, userService)), nil
}
