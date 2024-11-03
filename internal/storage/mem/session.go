package mem

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/shuryak/sberdevices/internal/model"
)

type SessionStorage struct {
	accessTokenLength  int
	refreshTokenLength int
	storage            sync.Map
}

func NewSessionStorage(accessTokenLength, refreshTokenLength int) *SessionStorage {
	return &SessionStorage{
		accessTokenLength:  accessTokenLength,
		refreshTokenLength: refreshTokenLength,
	}
}

func (s *SessionStorage) Create(
	_ context.Context,
	authCode, otpReceiverID, thirdPartyAccessToken, thirdPartyRefreshToken string,
	thirdPartyAccessTokenTTL, refreshTokenTTL time.Duration,
) (*model.Session, error) {
	session := model.NewSession(
		authCode,
		otpReceiverID,
		s.accessTokenLength,
		s.refreshTokenLength,
		thirdPartyAccessToken,
		thirdPartyRefreshToken,
		thirdPartyAccessTokenTTL,
		refreshTokenTTL,
	)

	s.storage.Store(session.AccessToken, session)

	return session, nil
}

func (s *SessionStorage) GetByAuthCode(_ context.Context, authCode string) (*model.Session, error) {
	var session *model.Session

	s.storage.Range(func(k, v interface{}) bool {
		vSession := v.(*model.Session)

		if vSession.AuthCode == authCode {
			session = vSession
			return false
		}

		return true
	})

	if session == nil {
		return nil, ErrSessionNotFound
	}

	return session, nil
}

func (s *SessionStorage) GetByAccessToken(_ context.Context, accessToken string) (*model.Session, error) {
	v, ok := s.storage.Load(accessToken)
	if !ok {
		return nil, ErrSessionNotFound
	}

	session := v.(*model.Session)

	expiresAt := session.UpdatedAt.Add(session.ThirdPartyAccessTokenTTL)

	if expiresAt.Before(time.Now()) {
		return nil, ErrAccessTokenExpired
	}

	return session, nil
}

func (s *SessionStorage) GetByRefreshToken(_ context.Context, refreshToken string) (*model.Session, error) {
	var session *model.Session

	s.storage.Range(func(k, v interface{}) bool {
		vSession := v.(*model.Session)

		if vSession.RefreshToken == refreshToken {
			session = vSession
			return false
		}

		return true
	})

	if session == nil {
		return nil, ErrSessionNotFound
	}

	return session, nil
}

func (s *SessionStorage) Refresh(
	_ context.Context,
	session *model.Session, thirdPartyAccessToken, thirdPartyRefreshToken string,
	thirdPartyTokenTTL, refreshTokenTTL time.Duration,
) (*model.Session, error) {
	var ns *model.Session

	ns = model.NewSession(
		session.AuthCode,
		session.OTPReceiverID,
		s.accessTokenLength,
		s.refreshTokenLength,
		thirdPartyAccessToken,
		thirdPartyRefreshToken,
		thirdPartyTokenTTL,
		refreshTokenTTL,
	)

	s.storage.Delete(session.AccessToken)
	s.storage.Store(ns.AccessToken, ns)

	return ns, nil
}

var (
	ErrAccessTokenExpired = errors.New("access token expired")
	ErrSessionNotFound    = errors.New("session not found")
)
