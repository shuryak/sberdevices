package adapter

import (
	"context"

	"github.com/shuryak/sberhack/internal/oauth"
	"github.com/shuryak/sberhack/pkg/pkce"
	"github.com/shuryak/sberhack/pkg/smarthome/auth"
)

type Authorizer struct {
	a *auth.Authorizer
}

func NewAuthorizer(a *auth.Authorizer) *Authorizer {
	return &Authorizer{a}
}

func (a *Authorizer) SendOTP(
	ctx context.Context,
	otpReceiverID string,
) (pkcePair *pkce.Pair, authOperationID string, err error) {
	return a.a.SendOTP(ctx, otpReceiverID)
}

func (a *Authorizer) GetTokensByOTP(
	ctx context.Context,
	authOperationID string,
	pkcePair *pkce.Pair,
	otp string,
) (*oauth.Tokens, error) {
	tokens, err := a.a.GetSmartHomeTokenByOTP(ctx, authOperationID, pkcePair, otp)
	if err != nil {
		return nil, err
	}

	return &oauth.Tokens{
		AccessToken:    tokens.SmartHomeToken,
		AccessTokenTTL: tokens.CSAFrontAccessTokenTTL,
		RefreshToken:   tokens.CSAFrontRefreshToken,
	}, err
}

func (a *Authorizer) RefreshTokens(ctx context.Context, refreshToken string) (*oauth.Tokens, error) {
	tokens, err := a.a.RefreshSmartHomeToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return &oauth.Tokens{
		AccessToken:    tokens.SmartHomeToken,
		AccessTokenTTL: tokens.CSAFrontAccessTokenTTL,
		RefreshToken:   tokens.CSAFrontRefreshToken,
	}, err
}
