package config_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pietro-putelli/feynman-backend/config"
)

var _ = Describe("Config", func() {
	Describe("NewConfig", func() {

		AfterEach(func() {
			os.Clearenv()
		})

		It("no env defined, should throw", func() {
			_, err := config.NewConfig()
			Expect(err).To(BeNil())
		})

		It("only databaes env defined, should throw", func() {
			os.Setenv("DB_HOST", "localhost")
			os.Setenv("DB_PORT", "5432")
			os.Setenv("DB_PASSWORD", "postgres")
			os.Setenv("DB_USER", "postgres")
			os.Setenv("DB_NAME", "postgres")

			_, err := config.NewConfig()
			Expect(err).To(HaveOccurred())
		})

		It("all required env defined, should succeed", func() {
			os.Setenv("DB_HOST", "localhost")
			os.Setenv("DB_PORT", "5432")
			os.Setenv("DB_PASSWORD", "postgres")
			os.Setenv("DB_USER", "postgres")
			os.Setenv("DB_NAME", "postgres")
			os.Setenv("AUTH_GOOGLE_CLIENT_ID", "google")
			os.Setenv("JWT_SECRET", "secret")

			cfg, err := config.NewConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
			Expect(cfg.Auth.Google.ClientID).To(Equal("google"))
			Expect(cfg.Auth.Jwt.Secret).To(Equal("secret"))
			Expect(cfg.Auth.Jwt.AccessTokenDuration).To(Equal(604800))
			Expect(cfg.Database.Host).To(Equal("localhost"))
			Expect(cfg.Database.Port).To(Equal(5432))
			Expect(cfg.Database.User).To(Equal("postgres"))
			Expect(cfg.Database.Password).To(Equal("postgres"))
			Expect(cfg.Database.Name).To(Equal("postgres"))
		})

		It("all env defined, should succeed", func() {
			os.Setenv("DB_HOST", "localhost")
			os.Setenv("DB_PORT", "5432")
			os.Setenv("DB_PASSWORD", "postgres")
			os.Setenv("DB_USER", "postgres")
			os.Setenv("DB_NAME", "postgres")
			os.Setenv("AUTH_GOOGLE_CLIENT_ID", "google")
			os.Setenv("JWT_SECRET", "secret")
			os.Setenv("JWT_ACCESS_TOKEN_DURATION", "3600")

			cfg, err := config.NewConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())
			Expect(cfg.Auth.Google.ClientID).To(Equal("google"))
			Expect(cfg.Auth.Jwt.Secret).To(Equal("secret"))
			Expect(cfg.Auth.Jwt.AccessTokenDuration).To(Equal(3600))
			Expect(cfg.Database.Host).To(Equal("localhost"))
			Expect(cfg.Database.Port).To(Equal(5432))
			Expect(cfg.Database.User).To(Equal("postgres"))
			Expect(cfg.Database.Password).To(Equal("postgres"))
			Expect(cfg.Database.Name).To(Equal("postgres"))
		})
	})
})
