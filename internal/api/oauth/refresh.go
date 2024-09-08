package oauth

import (
	"fmt"
	"net/http"

	"github.com/shuryak/sberhack/internal/api"
)

type RefreshResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type refreshReq struct {
	GrantType    string `query:"grant_type"`
	RefreshToken string `query:"refresh_token"`
	ClientID     string `query:"client_id"`
	ClientSecret string `query:"client_secret"`
}

func (p refreshReq) Validate(_ *api.Context) error {
	if p.GrantType != "refresh_token" {
		return fmt.Errorf("grant_type must be 'refresh_token'")
	}
	if len(p.RefreshToken) == 0 {
		return fmt.Errorf("refresh_token is required")
	}
	if len(p.ClientID) == 0 { // TODO: client_id from config
		return fmt.Errorf("invalid client_id")
	}
	if len(p.ClientSecret) == 0 { // TODO: client_secret from config
		return fmt.Errorf("invalid client_secret")
	}

	return nil
}

func (h *Handlers) Refresh(ctx *api.Context, req *refreshReq) (*RefreshResp, int) {
	h.log.Println("oauth: refresh")

	session, err := h.flow.RefreshSession(ctx, req.RefreshToken)
	if err != nil {
		h.log.Printf("refresh session failed, err: %v\n", err)
		_ = ctx.WriteResponse(http.StatusUnauthorized, nil)
		return nil, http.StatusUnauthorized // TODO: everywhere return error json object
	}

	return &RefreshResp{
		AccessToken:  session.AccessToken,
		TokenType:    "bearer",
		ExpiresIn:    int(session.ThirdPartyAccessTokenTTL.Seconds()),
		RefreshToken: session.RefreshToken,
	}, http.StatusOK
}
