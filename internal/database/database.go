package database

import (
	"fmt"

	"github.com/pietro-putelli/feynman-backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewConn creates a new database connection.
func NewConn(config *config.Database) gorm.Dialector {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		config.Host,
		config.User,
		config.Password,
		config.Name,
		config.Port,
		config.SSLMode,
	)

	conn := postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	})
	return conn
}

// NewConnection creates a new database connection.
func NewDB(dialector gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
