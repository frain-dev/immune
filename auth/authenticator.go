package auth

import (
	"net/http"

	"github.com/frain-dev/immune/config"
)

type Authenticator interface {
	Authenticate(r *http.Request) bool
}

func NewAuthenticator(cfg *config.EndpointConfig) Authenticator {
	if cfg.Authentication == nil {
		return noopAuthenticator{}
	}

	switch cfg.Authentication.Type {
	case "api_key":
		return &apiKeyAuthenticator{
			headerName:  cfg.Authentication.ApiKey.HeaderName,
			headerValue: cfg.Authentication.ApiKey.HeaderValue,
		}
	default:
		return noopAuthenticator{}
	}
}

type apiKeyAuthenticator struct {
	headerName  string
	headerValue string
}

func (a *apiKeyAuthenticator) Authenticate(r *http.Request) bool {
	val := r.Header.Get(a.headerName)
	return val == a.headerValue
}

type noopAuthenticator struct{}

func (a noopAuthenticator) Authenticate(r *http.Request) bool {
	return true
}
