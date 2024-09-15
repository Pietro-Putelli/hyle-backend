package main

import (
	"log"
	"os"

	"github.com/pietro-putelli/feynman-backend/openapi"
	"github.com/swaggest/openapi-go/openapi3"
)

type apiFunc func(openapi3.Reflector, string) error

const openApiVersion = "3.0.3"
const openApiTitle = "Faynman API"
const outputFileName = "docs/openapi.yaml"

// setServers sets the servers for the API
func setServers(reflector openapi3.Reflector) {
	devServerDescription := "Development server for Faynman API"
	reflector.Spec.Servers = append(reflector.Spec.Servers, openapi3.Server{
		URL:         "https://api.faynman.com",
		Description: &devServerDescription,
	})
	localServerDescription := "Local server for Faynman API"
	reflector.Spec.Servers = append(reflector.Spec.Servers, openapi3.Server{
		URL:         "http://localhost:3000",
		Description: &localServerDescription,
	})
}

// setSecurity sets the security for the API
func setSecurity(reflector openapi3.Reflector) string {
	securityName := "bearerAuth"
	reflector.
		SpecEns().
		SetHTTPBearerTokenSecurity(securityName, "JWT", "Bearer")
	return securityName
}

// setTags sets the tags for the API
func setTags(reflector openapi3.Reflector) {
	authTagDescription := "Authentication API"
	reflector.Spec.Tags = []openapi3.Tag{{Name: "Auth", Description: &authTagDescription}}
	booksTagDescription := "Books API"
	reflector.Spec.Tags = append(reflector.Spec.Tags, openapi3.Tag{Name: "Books", Description: &booksTagDescription})
	picksTagDescription := "Picks API"
	reflector.Spec.Tags = append(reflector.Spec.Tags, openapi3.Tag{Name: "Picks", Description: &picksTagDescription})
}

// exportSchema exports the schema to a file
func exportSchema(reflector openapi3.Reflector) {
	schema, err := reflector.Spec.MarshalYAML()
	if err != nil {
		log.Fatal(err)
	}

	// Print the schema to a file docs/openapi.yaml
	file, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	_, err = file.Write(schema)
}

func main() {
	reflector := openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{Openapi: openApiVersion}
	reflector.Spec.Info.
		WithTitle(openApiTitle).
		WithVersion("0.0.1").
		WithDescription("API documentation for Faynman API")

	setServers(reflector)
	securityName := setSecurity(reflector)
	setTags(reflector)

	apis := []apiFunc{
		openapi.BuildAuthAPI,  // Auth API
		openapi.BuildBooksAPI, // Books API
		openapi.BuildPicksAPI, // Picks API
	}

	for _, api := range apis {
		err := api(reflector, securityName)
		if err != nil {
			log.Fatal(err)
		}
	}

	exportSchema(reflector)
}
