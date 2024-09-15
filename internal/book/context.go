package book

import (
	"errors"

	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/pietro-putelli/feynman-backend/internal/database"
	"github.com/pietro-putelli/feynman-backend/internal/user"
	"gorm.io/gorm"
)

type Context struct {
	Service  Service
	Database *gorm.DB
}

func NewContext() (*Context, error) {
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

	userService := user.NewService(database)

	service := NewService(database, userService)

	return &Context{
		Service:  service,
		Database: database,
	}, nil
}
