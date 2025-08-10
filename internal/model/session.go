package model

import (
	"time"
)

type Session struct {
	OTPReceiverID string

	AccessToken  string
	RefreshToken string

	SmartHomeAccessToken    string
	SmartHomeAccessTokenTTL time.Duration

	SmartHomeRefreshToken string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSession(
	otpReceiverID, accessToken, refreshToken string, smartHomeAccessToken string,
	smartHomeAccessTokenTTL time.Duration, smartHomeRefreshToken string,
) *Session {
	now := time.Now()

	return &Session{
		OTPReceiverID:           otpReceiverID,
		AccessToken:             accessToken,
		RefreshToken:            refreshToken,
		SmartHomeAccessToken:    smartHomeAccessToken,
		SmartHomeAccessTokenTTL: smartHomeAccessTokenTTL,
		SmartHomeRefreshToken:   smartHomeRefreshToken,
		CreatedAt:               now,
		UpdatedAt:               now,
	}
}

func (s *Session) Refresh(
	accessToken, refreshToken, smartHomeAccessToken string, smartHomeAccessTokenTTL time.Duration,
	smartHomeRefreshToken string,
) *Session {
	return &Session{
		OTPReceiverID:           s.OTPReceiverID,
		AccessToken:             accessToken,
		RefreshToken:            refreshToken,
		SmartHomeAccessToken:    smartHomeAccessToken,
		SmartHomeAccessTokenTTL: smartHomeAccessTokenTTL,
		SmartHomeRefreshToken:   smartHomeRefreshToken,
		CreatedAt:               s.CreatedAt,
		UpdatedAt:               time.Now(),
	}
}
