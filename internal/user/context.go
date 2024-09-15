package user

import (
	"errors"

	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/pietro-putelli/feynman-backend/internal/database"
	"gorm.io/gorm"
)

type Context struct {
	Service  Service
	Database *gorm.DB
	Config   *config.Config
}

func NewContext() (*Context, error) {
	config, err := config.NewConfig()
	if err != nil {
		return nil, errors.New("failed load auth context config: " + err.Error())
	}

	database, err := database.NewDB(database.NewConn(&config.Database))
	if err != nil {
		return nil, errors.New("failed load auth context database: " + err.Error())
	}

	userService := NewService(database)

	return &Context{
		Service:  userService,
		Database: database,
		Config:   config,
	}, nil
}
