package database_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"

	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/pietro-putelli/feynman-backend/internal/database"
)

var _ = Describe("Database", func() {
	Describe("NewConn", func() {
		It("should return a new database connection", func() {
			// Arrange
			config := &config.Database{
				Host:     "localhost",
				User:     "user",
				Password: "password",
				Name:     "name",
				Port:     5432,
				SSLMode:  "disable",
			}

			// Act
			conn := database.NewConn(config)

			// Assert
			Expect(conn).NotTo(BeNil())
			Expect(conn).To(BeAssignableToTypeOf(&postgres.Dialector{}))
			Expect(conn.(*postgres.Dialector).Config.DSN).To(Equal("host=localhost user=user password=password dbname=name port=5432 sslmode=disable TimeZone=UTC"))
		})
	})

	Describe("NewDB", func() {
		It("should return a new database connection", func() {
			// Arrange
			mockDb, _, _ := sqlmock.New()
			conn := postgres.New(postgres.Config{
				Conn:       mockDb,
				DriverName: "postgres",
			})

			// Act
			db, err := database.NewDB(conn)

			// Assert
			Expect(err).To(BeNil())
			Expect(db).NotTo(BeNil())
		})

		It("should return an error when the connection fails", func() {
			// Arrange
			conn := postgres.New(postgres.Config{
				Conn:       nil,
				DriverName: "postgres",
			})

			// Act
			db, err := database.NewDB(conn)

			// Assert
			Expect(err).NotTo(BeNil())
			Expect(db).To(BeNil())
		})
	})
})
