package auth_test

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/pietro-putelli/feynman-backend/internal/auth"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/pietro-putelli/feynman-backend/internal/user"
)

var _ = Describe("Auth Service", func() {
	Describe("CreateUserIfNotExists", func() {
		var (
			service auth.Service
			// sqlMock sqlmock.Sqlmock
			userService *user.MockService
			ctrl        *gomock.Controller
		)

		authConfig := &config.AuthJwt{
			Secret:              "secret",
			AccessTokenDuration: 3600,
		}

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			userService = user.NewMockService(ctrl)
			// database, _ := database.NewDB(conn)
			service = auth.NewService(authConfig, userService)
		})

		It("should return the new user", func() {
			// Arrange
			user := &domain.ThirdPartyUser{
				Email:      "test@gmail.com",
				Sub:        "test123",
				GivenName:  "Given",
				FamilyName: "Family",
			}

			userResult := &domain.User{
				Guid:       uuid.Must(uuid.NewRandom()),
				ExternalID: user.Sub,
				Email:      user.Email,
				GivenName:  user.GivenName,
				FamilyName: user.FamilyName,
				Settings:   &domain.UserSettings{},
			}

			userService.EXPECT().CreateUserIfNotExists(user).Return(userResult, nil)

			// Act
			createdUser, err := service.CreateUserIfNotExists(user)

			// Assert
			Expect(err).To(BeNil())
			Expect(createdUser).NotTo(BeNil())
			Expect(createdUser.Email).To(Equal(user.Email))
			Expect(createdUser.ExternalID).To(Equal(user.Sub))
			Expect(createdUser.GivenName).To(Equal(user.GivenName))
			Expect(createdUser.FamilyName).To(Equal(user.FamilyName))
			Expect(createdUser.Settings).NotTo(BeNil())
			Expect(createdUser.Guid).ToNot(BeEmpty())
		})
	})

	Describe("GenerateToken", func() {
		var (
			service     auth.Service
			userService *user.MockService
		)

		authConfig := &config.AuthJwt{
			Secret:              "secret",
			AccessTokenDuration: 3600,
		}

		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			userService = user.NewMockService(ctrl)
			service = auth.NewService(authConfig, userService)
		})

		It("should generate a new token", func() {
			// Arrange
			user := &domain.User{
				Guid:       uuid.Must(uuid.NewRandom()),
				Email:      "test@gmail.com",
				GivenName:  "Given",
				FamilyName: "Family",
			}

			// Act
			token, err := service.GenerateToken(user)

			// Assert
			Expect(err).To(BeNil())
			Expect(token).NotTo(BeNil())
			Expect(token.AccessToken).NotTo(BeEmpty())
			Expect(token.RefreshToken).NotTo(BeEmpty())
			Expect(token.AccessTokenExpiresAt).NotTo(BeNil())
			// Parse the access token to check the claims
			jwtToken, err := jwt.Parse(token.AccessToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(authConfig.Secret), nil
			})
			Expect(err).To(BeNil())
			claims, ok := jwtToken.Claims.(jwt.MapClaims)
			Expect(ok).To(BeTrue())
			Expect(claims["sub"]).To(Equal(user.Guid.String()))
			Expect(claims["email"]).To(Equal(user.Email))
			Expect(claims["iss"]).To(Equal("panta.srvless-api"))
			Expect(claims["kind"]).To(Equal("access"))

			// Parse the refresh token to check the claims
			jwtToken, err = jwt.Parse(token.RefreshToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(authConfig.Secret), nil
			})
			Expect(err).To(BeNil())
			claims, ok = jwtToken.Claims.(jwt.MapClaims)
			Expect(ok).To(BeTrue())
			Expect(claims["sub"]).To(Equal(user.Guid.String()))
			Expect(claims["email"]).To(Equal(user.Email))
			Expect(claims["iss"]).To(Equal("panta.srvless-api"))
			Expect(claims["kind"]).To(Equal("refresh"))
		})
	})

	Describe("RefreshAccessToken", func() {
		var (
			service     auth.Service
			userService *user.MockService
		)

		authConfig := &config.AuthJwt{
			Secret:              "secret",
			AccessTokenDuration: 3600,
		}

		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			userService = user.NewMockService(ctrl)
			service = auth.NewService(authConfig, userService)
		})

		It("should generate a new access token", func() {
			// Arrange
			user := &domain.User{
				Guid:       uuid.Must(uuid.NewRandom()),
				Email:      "test@test,com",
				GivenName:  "Given",
				FamilyName: "Family",
			}

			// Generate a refresh token
			refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub":   user.Guid,
				"email": user.Email,
				"exp":   time.Now().AddDate(1000, 0, 0).Unix(),
				"iat":   time.Now().Unix(),
				"iss":   "panta.srvless-api",
				"kind":  "refresh",
			})
			refreshTokenString, err := refreshToken.SignedString([]byte(authConfig.Secret))
			Expect(err).To(BeNil())

			userService.EXPECT().GetUserByGuid(user.Guid).Return(user, nil)

			// Act
			token, err := service.RefreshAccessToken(refreshTokenString)

			// Assert
			Expect(err).To(BeNil())
			Expect(token).NotTo(BeNil())
			Expect(token.AccessToken).NotTo(BeEmpty())
			Expect(token.RefreshToken).NotTo(BeEmpty())
			Expect(token.AccessTokenExpiresAt).NotTo(BeNil())
			// Parse the access token to check the claims
			jwtToken, err := jwt.Parse(token.AccessToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(authConfig.Secret), nil
			})
			Expect(err).To(BeNil())
			claims, ok := jwtToken.Claims.(jwt.MapClaims)
			Expect(ok).To(BeTrue())
			Expect(claims["sub"]).To(Equal(user.Guid.String()))
			Expect(claims["email"]).To(Equal(user.Email))
			Expect(claims["iss"]).To(Equal("panta.srvless-api"))
			Expect(claims["kind"]).To(Equal("access"))
		})
	})
})
