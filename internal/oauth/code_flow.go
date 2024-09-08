package oauth

import (
	"context"
	"fmt"
	"time"

	"github.com/shuryak/sberhack/internal/model"
	"github.com/shuryak/sberhack/pkg/pkce"
)

type CodeStorage interface {
	Create(
		ctx context.Context,
		thirdPartyPKCEPair *pkce.Pair,
		thirdPartyOTPReceiverID, thirdPartyAuthOperationID string,
	) (authCode string, err error)
	Get(ctx context.Context, authCode string) (*model.AuthCodePayload, error)
	Delete(ctx context.Context, authCode string) error
}

// ThirdPartyAuthProvider
// TODO: describe otpReceiverID == phoneNumber, authOperationID == OUID
// TODO: encapsulate pkcePair on new struct type crossRequestSpecificData
type ThirdPartyAuthProvider interface {
	SendOTP(ctx context.Context, otpReceiverID string) (pkcePair *pkce.Pair, authOperationID string, err error)
	GetTokensByOTP(ctx context.Context, authOperationID string, pkcePair *pkce.Pair, otp string) (*Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*Tokens, error)
}

type SessionStorage interface {
	Create(
		ctx context.Context,
		authCode, otpReceiverID, thirdPartyAccessToken, thirdPartyRefreshToken string,
		thirdPartyTokenTTL, refreshTokenTTL time.Duration,
	) (*model.Session, error)
	GetByAuthCode(ctx context.Context, authCode string) (*model.Session, error)
	GetByAccessToken(ctx context.Context, accessToken string) (*model.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error)
	Refresh(
		ctx context.Context,
		session *model.Session, thirdPartyAccessToken, thirdPartyRefreshToken string,
		thirdPartyTokenTTL, refreshTokenTTL time.Duration,
	) (*model.Session, error)
}

type Tokens struct {
	AccessToken    string
	AccessTokenTTL time.Duration
	RefreshToken   string
}

type CodeFlowWithOTP struct {
	codeStorage    CodeStorage
	sessionStorage SessionStorage
	authProvider   ThirdPartyAuthProvider
}

func NewCodeFlowWithOTP(cs CodeStorage, ss SessionStorage, ap ThirdPartyAuthProvider) *CodeFlowWithOTP {
	return &CodeFlowWithOTP{
		codeStorage:    cs,
		sessionStorage: ss,
		authProvider:   ap,
	}
}

func (f *CodeFlowWithOTP) Start(ctx context.Context, otpReceiverID string) (authCode string, err error) {
	pkcePair, authOperationID, err := f.authProvider.SendOTP(ctx, otpReceiverID)
	if err != nil {
		return "", fmt.Errorf("failed to send otp, err: %w", err)
	}

	return f.codeStorage.Create(ctx, pkcePair, otpReceiverID, authOperationID)
}

func (f *CodeFlowWithOTP) Get(ctx context.Context, authCode string) (*model.AuthCodePayload, error) {
	return f.codeStorage.Get(ctx, authCode)
}

func (f *CodeFlowWithOTP) Delete(ctx context.Context, authCode string) error {
	return f.codeStorage.Delete(ctx, authCode)
}

func (f *CodeFlowWithOTP) CreateSession(ctx context.Context, authCode, otp string) (*model.Session, error) {
	// TODO: check for access_token existing

	codePayload, err := f.codeStorage.Get(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth code payload, err: %w", err)
	}

	err = f.codeStorage.Delete(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to delete auth code, err: %w", err)
	}

	tokens, err := f.authProvider.GetTokensByOTP(
		ctx,
		codePayload.ThirdPartyAuthOperationID,
		codePayload.ThirdPartyPKCEPair,
		otp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange otp, err: %w", err)
	}

	return f.sessionStorage.Create(
		context.Background(),
		authCode,
		codePayload.ThirdPartyOTPReceiverID,
		tokens.AccessToken,
		tokens.RefreshToken,
		tokens.AccessTokenTTL,
		time.Hour*24*30, // TODO: ttl
	)
}

func (f *CodeFlowWithOTP) GetSessionByAuthCode(ctx context.Context, authCode string) (*model.Session, error) {
	return f.sessionStorage.GetByAuthCode(ctx, authCode)
}

func (f *CodeFlowWithOTP) GetSessionByAccessToken(ctx context.Context, accessToken string) (*model.Session, error) {
	return f.sessionStorage.GetByAccessToken(ctx, accessToken)
}

func (f *CodeFlowWithOTP) RefreshSession(ctx context.Context, refreshToken string) (*model.Session, error) {
	session, err := f.sessionStorage.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	tokens, err := f.authProvider.RefreshTokens(ctx, session.ThirdPartyRefreshToken)
	if err != nil {
		return nil, err
	}

	return f.sessionStorage.Refresh(
		context.Background(),
		session,
		tokens.AccessToken,
		tokens.RefreshToken,
		tokens.AccessTokenTTL,
		time.Hour*24*30, // TODO: ttl
	)
}
