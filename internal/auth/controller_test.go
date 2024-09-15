package auth_test

import (
	"errors"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/pietro-putelli/feynman-backend/internal/auth"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
)

var _ = Describe("Auth Controller", func() {
	var (
		ctrl *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("HandleAuthToken", func() {
		cfg := &config.Config{
			Auth: config.Auth{
				Jwt: config.AuthJwt{
					Secret: "jwt-secret",
				},
			},
		}

		Context("when provided auth provider is not valid", func() {
			It("should return an error", func() {
				// Arrange
				body := &domain.AuthTokenBody{
					Token:    "example-token",
					Provider: "not-supported-provider",
				}

				providerHandler := func(provider string, config *config.Auth) (auth.ThirdPartyProvider, error) {
					return nil, errors.New("provider not supported")
				}

				ctx := auth.NewContext(cfg, nil, providerHandler, nil)

				// Act
				result, err := auth.HandleAuthToken(ctx, body)

				// Assert
				Expect(result).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("provider not supported"))
			})
		})

		Context("when provided token is not valid", func() {
			It("should return an error", func() {
				// Arrange
				body := &domain.AuthTokenBody{
					Token:    "example-token",
					Provider: "example-provider",
				}

				mockThirdPartyProvider := auth.NewMockThirdPartyProvider(ctrl)
				mockThirdPartyProvider.EXPECT().ValidateToken(body.Token).Return(nil, errors.New("invalid token"))

				providerHandler := func(provider string, config *config.Auth) (auth.ThirdPartyProvider, error) {
					return mockThirdPartyProvider, nil
				}

				ctx := auth.NewContext(cfg, nil, providerHandler, nil)

				// Act
				result, err := auth.HandleAuthToken(ctx, body)

				// Assert
				Expect(result).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("invalid token"))
			})
		})

		Context("when token is valid but service fails to handle user", func() {
			It("should return an error", func() {
				// Arrange
				body := &domain.AuthTokenBody{
					Token:    "example-token",
					Provider: "example-provider",
				}
				thirdPartyUser := &domain.ThirdPartyUser{
					Email:      "test@gmail.com",
					Sub:        "123",
					Name:       "Test User",
					GivenName:  "Test",
					FamilyName: "User",
				}

				mockThirdPartyProvider := auth.NewMockThirdPartyProvider(ctrl)
				mockThirdPartyProvider.EXPECT().ValidateToken(body.Token).Return(thirdPartyUser, nil)

				mockService := auth.NewMockService(ctrl)
				mockService.EXPECT().CreateUserIfNotExists(thirdPartyUser).Return(nil, errors.New("failed to create user"))

				providerHandler := func(provider string, config *config.Auth) (auth.ThirdPartyProvider, error) {
					return mockThirdPartyProvider, nil
				}

				ctx := auth.NewContext(cfg, nil, providerHandler, mockService)

				// Act
				result, err := auth.HandleAuthToken(ctx, body)

				// Assert
				Expect(result).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to create user"))
			})
		})

		Context("when token is valid, user is created and service fails to generate token", func() {
			It("should return an error", func() {
				// Arrange
				body := &domain.AuthTokenBody{
					Token:    "example-token",
					Provider: "example-provider",
				}
				thirdPartyUser := &domain.ThirdPartyUser{
					Email:      "test@gmail.com",
					Sub:        "123",
					Name:       "Test User",
					GivenName:  "Test",
					FamilyName: "User",
				}
				user := &domain.User{
					ID:         0,
					Guid:       uuid.Must(uuid.NewRandom()),
					Email:      thirdPartyUser.Email,
					ExternalID: thirdPartyUser.Sub,
					GivenName:  thirdPartyUser.GivenName,
					Settings:   domain.NewUserSettings(),
					FamilyName: thirdPartyUser.FamilyName,
				}

				mockThirdPartyProvider := auth.NewMockThirdPartyProvider(ctrl)
				mockThirdPartyProvider.EXPECT().ValidateToken(body.Token).Return(thirdPartyUser, nil)

				mockService := auth.NewMockService(ctrl)
				mockService.EXPECT().CreateUserIfNotExists(thirdPartyUser).Return(user, nil)
				mockService.EXPECT().GenerateToken(user).Return(nil, errors.New("failed to generate token"))

				providerHandler := func(provider string, config *config.Auth) (auth.ThirdPartyProvider, error) {
					return mockThirdPartyProvider, nil
				}

				ctx := auth.NewContext(cfg, nil, providerHandler, mockService)

				// Act
				result, err := auth.HandleAuthToken(ctx, body)

				// Assert
				Expect(result).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to generate token"))
			})
		})

		Context("when token is valid, user is created and token generated", func() {
			It("should return the user and token", func() {
				// Arrange
				body := &domain.AuthTokenBody{
					Token:    "example-token",
					Provider: "example-provider",
				}
				thirdPartyUser := &domain.ThirdPartyUser{
					Email:      "test@gmail.com",
					Sub:        "123",
					Name:       "Test User",
					GivenName:  "Test",
					FamilyName: "User",
				}
				user := &domain.User{
					ID:         0,
					Guid:       uuid.Must(uuid.NewRandom()),
					Email:      thirdPartyUser.Email,
					ExternalID: thirdPartyUser.Sub,
					GivenName:  thirdPartyUser.GivenName,
					Settings:   domain.NewUserSettings(),
					FamilyName: thirdPartyUser.FamilyName,
				}
				authToken := &domain.AuthTokenDto{
					AccessToken:          "access-token",
					RefreshToken:         "refresh-token",
					AccessTokenExpiresAt: time.Now().Add(time.Hour),
				}

				mockThirdPartyProvider := auth.NewMockThirdPartyProvider(ctrl)
				mockThirdPartyProvider.EXPECT().ValidateToken(body.Token).Return(thirdPartyUser, nil)

				mockService := auth.NewMockService(ctrl)
				mockService.EXPECT().CreateUserIfNotExists(thirdPartyUser).Return(user, nil)
				mockService.EXPECT().GenerateToken(user).Return(authToken, nil)

				providerHandler := func(provider string, config *config.Auth) (auth.ThirdPartyProvider, error) {
					return mockThirdPartyProvider, nil
				}

				ctx := auth.NewContext(cfg, nil, providerHandler, mockService)

				// Act
				result, err := auth.HandleAuthToken(ctx, body)

				// Assert
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.User).NotTo(BeNil())
				Expect(result.User.Email).To(Equal(user.Email))
				Expect(result.User.GivenName).To(Equal(user.GivenName))
				Expect(result.User.FamilyName).To(Equal(user.FamilyName))
				Expect(result.Auth).NotTo(BeNil())
				Expect(result.Auth.AccessToken).To(Equal(authToken.AccessToken))
				Expect(result.Auth.RefreshToken).To(Equal(authToken.RefreshToken))
				Expect(result.Auth.AccessTokenExpiresAt).To(Equal(authToken.AccessTokenExpiresAt))
			})
		})
	})
})
