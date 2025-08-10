package oauth

import (
	"errors"
	"net/http"

	"github.com/shuryak/sberdevices/internal/api"
)

type ErrorResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func (h *Handlers) ErrHandler(_ *api.Context, err error) (interface{}, int) {
	resp := ErrorResp{}
	statusCode := http.StatusBadRequest

	switch {
	case errors.Is(err, ErrInvalidRedirectURI):
		resp.Error = ErrOAuthInvalidRequest.Error()
		resp.ErrorDescription = "redirect_uri is invalid"
	case errors.Is(err, ErrOAuthUnsupportedGrantType):
		resp.Error = err.Error()
		resp.ErrorDescription = "grant_type must be 'authorization_code'"
	default:
		resp.Error = ErrOAuthInvalidRequest.Error()
	}

	return resp, statusCode
}

func IsOAuthError(err error) bool {
	return errors.Is(err, ErrOAuthInvalidRequest) ||
		errors.Is(err, ErrOAuthInvalidClient) ||
		errors.Is(err, ErrOAuthInvalidGrant) ||
		errors.Is(err, ErrOAuthUnauthorizedClient) ||
		errors.Is(err, ErrOAuthUnsupportedGrantType) ||
		errors.Is(err, ErrOAuthInvalidScope)
}

// https://datatracker.ietf.org/doc/html/rfc6749#section-5.2
var (
	ErrOAuthInvalidRequest       = errors.New("invalid_request")
	ErrOAuthInvalidClient        = errors.New("invalid_client")
	ErrOAuthInvalidGrant         = errors.New("invalid_grant")
	ErrOAuthUnauthorizedClient   = errors.New("unauthorized_client")
	ErrOAuthUnsupportedGrantType = errors.New("unsupported_grant_type")
	ErrOAuthInvalidScope         = errors.New("invalid_scope")
)
