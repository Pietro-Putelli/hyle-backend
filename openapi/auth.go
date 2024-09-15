package openapi

import (
	"net/http"

	"github.com/pietro-putelli/feynman-backend/internal/domain"
	"github.com/swaggest/openapi-go/openapi3"
)

// BuildAuthAPI builds the auth openapi API endpoints
func BuildAuthAPI(reflector openapi3.Reflector, securityName string) error {
	// AuthTokenPostFun
	authTokenOp, err := reflector.NewOperationContext(http.MethodPost, "/v1/auth/token")
	if err != nil {
		return err
	}
	authTokenOp.SetSummary("Retrieve auth token")
	authTokenOp.SetDescription("Retrieve the authentication token")
	authTokenOp.SetID("authGetToken")
	authTokenOp.AddRespStructure(domain.AuthTokenDto{})
	authTokenOp.AddReqStructure(domain.AuthTokenBody{})
	authTokenOp.SetTags("Auth")

	err = reflector.AddOperation(authTokenOp)
	if err != nil {
		return err
	}

	// AuthRefreshPostFun
	authRefreshOp, err := reflector.NewOperationContext(http.MethodPost, "/v1/auth/refresh")
	if err != nil {
		return err
	}
	authRefreshOp.AddSecurity(securityName)
	authRefreshOp.SetSummary("Refresh auth token")
	authRefreshOp.SetDescription("Refresh the authentication token, the bearer token must be set to the refresh token")
	authRefreshOp.SetID("authRefreshToken")
	authRefreshOp.AddRespStructure(domain.AuthTokenDto{})
	authRefreshOp.SetTags("Auth")

	err = reflector.AddOperation(authRefreshOp)
	if err != nil {
		return err
	}

	return nil
}
