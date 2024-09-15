package auth_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/pietro-putelli/feynman-backend/internal/auth"
)

var _ = Describe("Auth Provider", func() {

	Describe("NewProvider", func() {
		It("google, should return google provider", func() {
			// Arrange
			provider := "google"
			config := &config.Auth{
				Google: config.AuthGoogle{
					ClientID: "client-id",
				},
			}

			// Act
			result, err := auth.NewProvider(provider, config)

			// Assert
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result).To(BeAssignableToTypeOf(&auth.GoogleProvider{}))
		})

		It("apple, should return apple provider", func() {
			// Arrange
			provider := "apple"
			config := &config.Auth{}

			// Act
			result, err := auth.NewProvider(provider, config)

			// Assert
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result).To(BeAssignableToTypeOf(&auth.AppleProvider{}))
		})

		It("invalid provider, should return error", func() {
			// Arrange
			provider := "invalid"
			config := &config.Auth{}

			// Act
			result, err := auth.NewProvider(provider, config)

			// Assert
			Expect(err).NotTo(BeNil())
			Expect(result).To(BeNil())
		})
	})
})
