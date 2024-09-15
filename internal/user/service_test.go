package user_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/pietro-putelli/feynman-backend/internal/database"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/user"
)

var _ = Describe("Service", func() {
	Describe("NewService", func() {
		It("should return user service", func() {
			// Arrange
			// Act
			result := user.NewService(nil)

			// Assert
			Expect(result).NotTo(BeNil())
		})
	})

	Describe("CreateUserIfNotExists", func() {
		var (
			service user.Service
			sqlMock sqlmock.Sqlmock
		)

		BeforeEach(func() {
			db, sqlMockGen, _ := sqlmock.New()
			sqlMock = sqlMockGen

			conn := postgres.New(postgres.Config{
				Conn: db,
			})

			database, _ := database.NewDB(conn)
			service = user.NewService(database)
		})

		It("should create a new user", func() {
			// Arrange
			user := &domain.ThirdPartyUser{
				Email:      "test@test.com",
				Sub:        "test123",
				GivenName:  "Given",
				FamilyName: "Family",
				Provider:   "google",
			}

			expectedSelect := "^SELECT (.+) FROM \"users\" WHERE email = (.+)$"
			sqlMock.ExpectQuery(expectedSelect).WillReturnError(gorm.ErrRecordNotFound)
			sqlMock.ExpectBegin()
			sqlMock.ExpectQuery("^INSERT INTO \"users\" (.+) RETURNING (.+)$").
				WithArgs(user.Sub, user.Email, user.GivenName, user.FamilyName, user.Provider, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnRows(sqlMock.NewRows([]string{"id"}).AddRow(1))
			sqlMock.ExpectCommit()

			// Act
			result, sessionID, err := service.CreateUserIfNotExists(user)

			// Assert
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(sessionID).NotTo(BeEmpty())
			Expect(result.Email).To(Equal(user.Email))
			Expect(result.GivenName).To(Equal(user.GivenName))
			Expect(result.FamilyName).To(Equal(user.FamilyName))
		})

		It("should return an existing user", func() {
			// Arrange
			user := &domain.ThirdPartyUser{
				Email:      "test@test.com",
				Sub:        "test123",
				GivenName:  "Given",
				FamilyName: "Family",
			}

			// return user
			expectedSelect := "^SELECT (.+) FROM \"users\" WHERE email = (.+)$"
			sqlMock.ExpectQuery(expectedSelect).
				WithArgs(user.Email, 1).
				WillReturnRows(sqlMock.NewRows([]string{"id", "email", "given_name", "family_name", "external_id"}).
					AddRow(1, user.Email, user.GivenName, user.FamilyName, user.Sub))

			// Act
			result, sessionID, err := service.CreateUserIfNotExists(user)

			// Assert
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(sessionID).NotTo(BeEmpty())
			Expect(result.Email).To(Equal(user.Email))
			Expect(result.GivenName).To(Equal(user.GivenName))
			Expect(result.FamilyName).To(Equal(user.FamilyName))
			Expect(result.ExternalID).To(Equal(user.Sub))
		})

		It("select fails, should return an error", func() {
			// Arrange
			user := &domain.ThirdPartyUser{
				Email:      "test@test.com",
				Sub:        "test123",
				GivenName:  "Given",
				FamilyName: "Family",
			}

			expectedSelect := "^SELECT (.+) FROM \"users\" WHERE email = (.+)$"
			sqlMock.ExpectQuery(expectedSelect).WillReturnError(gorm.ErrInvalidDB)

			// Act
			result, sessionID, err := service.CreateUserIfNotExists(user)

			// Assert
			Expect(err).ToNot(BeNil())
			Expect(result).To(BeNil())
			Expect(sessionID).To(BeEmpty())
		})

		It("insert fails, should return an error", func() {
			// Arrange
			user := &domain.ThirdPartyUser{
				Email: "test@test.com",
			}

			expectedSelect := "^SELECT (.+) FROM \"users\" WHERE email = (.+)$"
			sqlMock.ExpectQuery(expectedSelect).WillReturnError(gorm.ErrRecordNotFound)
			sqlMock.ExpectBegin()
			sqlMock.ExpectQuery("^INSERT INTO \"users\" (.+) RETURNING").WillReturnError(gorm.ErrInvalidDB)
			sqlMock.ExpectRollback()

			// Act
			result, sessionID, err := service.CreateUserIfNotExists(user)

			// Assert
			Expect(err).ToNot(BeNil())
			Expect(result).To(BeNil())
			Expect(sessionID).To(BeEmpty())
		})
	})

	Describe("GetUserByGuid", func() {
		var (
			service user.Service
			sqlMock sqlmock.Sqlmock
		)

		BeforeEach(func() {
			db, sqlMockGen, _ := sqlmock.New()
			sqlMock = sqlMockGen

			conn := postgres.New(postgres.Config{
				Conn: db,
			})

			database, _ := database.NewDB(conn)
			service = user.NewService(database)
		})

		It("should return a user", func() {
			// Arrange
			guid := uuid.Must(uuid.NewRandom())

			user := &domain.User{
				Guid:       guid,
				ExternalID: "test123",
				Email:      "test@test.com",
				GivenName:  "Given",
				FamilyName: "Family",
			}

			expectedSelect := "^SELECT (.+) FROM \"users\" WHERE guid = (.+)$"
			sqlMock.ExpectQuery(expectedSelect).
				WithArgs(guid, 1).
				WillReturnRows(sqlMock.NewRows([]string{"id", "email", "given_name", "family_name", "external_id"}).
					AddRow(1, user.Email, user.GivenName, user.FamilyName, user.ExternalID))

			// Act
			result, err := service.GetUserByGuid(guid)

			// Assert
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Email).To(Equal(user.Email))
			Expect(result.GivenName).To(Equal(user.GivenName))
			Expect(result.FamilyName).To(Equal(user.FamilyName))
			Expect(result.ExternalID).To(Equal(user.ExternalID))
		})

		It("should return an error", func() {
			// Arrange
			guid := uuid.Must(uuid.NewRandom())

			expectedSelect := "^SELECT (.+) FROM \"users\" WHERE guid = (.+)$"
			sqlMock.ExpectQuery(expectedSelect).WillReturnError(gorm.ErrInvalidDB)

			// Act
			result, err := service.GetUserByGuid(guid)

			// Assert
			Expect(err).ToNot(BeNil())
			Expect(result).To(BeNil())
		})
	})
})
