package model

import (
	"time"

	"github.com/shuryak/sberhack/pkg/strrand"
)

type Session struct {
	AuthCode                 string
	OTPReceiverID            string
	AccessToken              string
	RefreshToken             string
	ThirdPartyAccessToken    string
	ThirdPartyRefreshToken   string
	ThirdPartyAccessTokenTTL time.Duration
	RefreshTokenExpiresTTL   time.Duration
	UpdatedAt                time.Time
}

func NewSession(
	authCode string,
	otpReceiverID string,
	accessTokenLength, refreshTokenLength int,
	thirdPartyAccessToken, thirdPartyRefreshToken string,
	thirdPartyTokenTTL, refreshTokenTTL time.Duration,
) *Session {
	return &Session{
		AuthCode:                 authCode,
		OTPReceiverID:            otpReceiverID,
		AccessToken:              strrand.RandSeqStr(accessTokenLength),
		RefreshToken:             strrand.RandSeqStr(refreshTokenLength),
		ThirdPartyAccessToken:    thirdPartyAccessToken,
		ThirdPartyRefreshToken:   thirdPartyRefreshToken,
		ThirdPartyAccessTokenTTL: thirdPartyTokenTTL,
		RefreshTokenExpiresTTL:   refreshTokenTTL,
		UpdatedAt:                time.Now(),
	}
}
