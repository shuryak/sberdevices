package oauth

import (
	"errors"
	"net/http"
	"net/url"
	"time"
	"unicode/utf8"

	"github.com/shuryak/sberdevices/internal/api"
	"github.com/shuryak/sberdevices/internal/config"
)

type tokenReq struct {
	GrantType    string `query:"grant_type"`
	Code         string `query:"code"`
	RedirectURI  string `query:"redirect_uri"`
	ClientID     string `query:"client_id"`
	ClientSecret string `query:"client_secret"`
}

func (p tokenReq) Validate(_ *api.Context) error {
	if p.GrantType != "authorization_code" {
		return ErrOAuthUnsupportedGrantType
	}
	if _, err := url.ParseRequestURI(p.RedirectURI); err != nil {
		return ErrInvalidRedirectURI
	}

	return nil
}

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (p tokenReq) AdditionalValidate(_ *api.Context, cfg *config.Config) error {
	if uint(utf8.RuneCountInString(p.Code)) != cfg.Auth.CodeLength {
		return errors.New("invalid code")
	}
	// if cfg.Clients.CheckClient(p.ClientID, p.ClientSecret) {
	// 	return errors.New("invalid client_id and client_secret pair")
	// }

	return nil
}

// Token - https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.3
func (h *Handlers) Token(ctx *api.Context, req *tokenReq) (*TokenResp, int) {
	if err := req.AdditionalValidate(ctx, h.cfg); err != nil {
		if !ctx.InvokeErrHandler(err) {
			h.log.Println("no error handler provided!")
		}
		return nil, 0
	}

	h.log.Println("oauth: token")

	session, err := h.flow.ExchangeAuthCode(ctx, req.Code)
	if err != nil {
		h.log.Printf("get session by auth code failed, err: %v\n", err)
		return nil, http.StatusBadRequest
	}

	return &TokenResp{
		AccessToken: session.AccessToken,
		TokenType:   "bearer",
		// ExpiresIn:    int(session.SmartHomeAccessTokenTTL.Seconds()),
		ExpiresIn:    int((15 * time.Second).Seconds()),
		RefreshToken: session.RefreshToken,
	}, http.StatusOK
}

var ErrInvalidRedirectURI = errors.New("invalid redirect_uri")
