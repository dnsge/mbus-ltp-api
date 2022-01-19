package mbus

import (
	"net/http"
)

type AuthApplier interface {
	ApplyAuth(req *http.Request) error
}

type APIKeyAuth struct {
	key string
}

func NewAPIKeyAuth(key string) *APIKeyAuth {
	return &APIKeyAuth{
		key: key,
	}
}

func (a *APIKeyAuth) ApplyAuth(req *http.Request) error {
	args := req.URL.Query()
	args.Set("key", a.key)
	req.URL.RawQuery = args.Encode()
	return nil
}
