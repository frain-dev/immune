package auth

import (
	"net/http"

	"github.com/frain-dev/immune/config"
)

type Authenticator interface {
	Authenticate(r *http.Request) bool
}

func NewAuthenticator(cfg *config.Config) Authenticator {
	if cfg.EndpointAuthentication == nil {
		return noopAuthenticator{}
	}

	switch cfg.EndpointAuthentication.Type {
	case "api_key":
		return &apiKeyAuthenticator{
			headerName:  cfg.EndpointAuthentication.ApiKey.HeaderName,
			headerValue: cfg.EndpointAuthentication.ApiKey.HeaderValue,
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
