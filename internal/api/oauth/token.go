package oauth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/shuryak/sberhack/internal/api"
)

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type tokenReq struct {
	GrantType    string `query:"grant_type"`
	Code         string `query:"code"`
	RedirectURI  string `query:"redirect_uri"`
	ClientID     string `query:"client_id"`
	ClientSecret string `query:"client_secret"`
}

func (p tokenReq) Validate(_ *api.Context) error {
	if p.GrantType != "authorization_code" {
		return fmt.Errorf("grant_type must be 'authorization_code'")
	}
	// if utf8.RuneCountInString(p.Code) != 64 { // TODO: code length from config
	// 	return fmt.Errorf("invalid code")
	// }
	if _, err := url.ParseRequestURI(p.RedirectURI); err != nil {
		return fmt.Errorf("invalid redirect_uri")
	}
	if len(p.ClientID) == 0 { // TODO: client_id from config
		return fmt.Errorf("invalid client_id")
	}
	if len(p.ClientSecret) == 0 { // TODO: client_secret from config
		return fmt.Errorf("invalid client_secret")
	}

	return nil
}

func (h *Handlers) Token(ctx *api.Context, req *tokenReq) (*TokenResp, int) {
	h.log.Println("oauth: token")

	session, err := h.flow.GetSessionByAuthCode(ctx, req.Code)
	if err != nil {
		h.log.Printf("get session by auth code failed, err: %v\n", err)
		return nil, http.StatusBadRequest
	}

	return &TokenResp{
		AccessToken:  session.AccessToken,
		TokenType:    "bearer",
		ExpiresIn:    int(session.ThirdPartyAccessTokenTTL.Seconds()),
		RefreshToken: session.RefreshToken,
	}, http.StatusOK
}
