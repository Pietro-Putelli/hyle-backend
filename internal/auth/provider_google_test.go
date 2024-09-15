package auth_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pietro-putelli/feynman-backend/internal/auth"
)

var _ = Describe("ProviderGoogle", func() {

	Describe("NewGoogleProvider", func() {
		It("should return google provider", func() {
			// Arrange
			clientID := "client-id"

			// Act
			result, err := auth.NewGoogleProvider(clientID)

			// Assert
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result).To(BeAssignableToTypeOf(&auth.GoogleProvider{}))
		})
	})
})
