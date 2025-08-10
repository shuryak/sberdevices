package oauth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shuryak/sberdevices/internal/model"
	"github.com/shuryak/sberdevices/internal/pkg/pkce"
	"github.com/shuryak/sberdevices/internal/pkg/strrand"
)

type CodeStorage interface {
	Create(
		ctx context.Context,
		thirdPartyPKCEPair *pkce.Pair,
		thirdPartyOTPReceiverID, thirdPartyAuthOperationID string,
	) (authCode string, err error)
	SetAccessToken(ctx context.Context, authCode, accessToken string) error
	Get(ctx context.Context, authCode string) (*model.AuthCodePayload, error)
	DeleteCodeAndGetSession(ctx context.Context, authCode string) (*model.Session, error)
}

// SmartHomeAuthProvider
// TODO: describe otpReceiverID == phoneNumber, authOperationID == OUID
// TODO: encapsulate pkcePair on new struct type crossRequestSpecificData
type SmartHomeAuthProvider interface {
	SendOTP(ctx context.Context, otpReceiverID string) (pkcePair *pkce.Pair, authOperationID string, err error)
	GetTokensByOTP(ctx context.Context, authOperationID string, pkcePair *pkce.Pair, otp string) (*Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*Tokens, error)
}

type SessionStorage interface {
	Upsert(ctx context.Context, oldAccessToken string, session *model.Session) (*model.Session, error)
	GetByAccessToken(ctx context.Context, accessToken string) (*model.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error)
}

type StorageUnitOfWork interface {
	SessionStorage() SessionStorage
	OAuthCodeStorage() CodeStorage
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type StorageUnitOfWorkFactory func(ctx context.Context) (StorageUnitOfWork, error)

type Tokens struct {
	AccessToken    string
	AccessTokenTTL time.Duration
	RefreshToken   string
}

type CodeFlowWithOTP struct {
	codeStorage       CodeStorage
	sessionStorage    SessionStorage
	storageUoWFactory StorageUnitOfWorkFactory
	authProvider      SmartHomeAuthProvider

	accessTokenLength  uint
	refreshTokenLength uint
}

func NewCodeFlowWithOTP(
	codeStorage CodeStorage, sessionStorage SessionStorage, uowFactory StorageUnitOfWorkFactory,
	authProvider SmartHomeAuthProvider, accessTokenLength, refreshTokenLength uint,
) *CodeFlowWithOTP {
	return &CodeFlowWithOTP{
		codeStorage:        codeStorage,
		sessionStorage:     sessionStorage,
		storageUoWFactory:  uowFactory,
		authProvider:       authProvider,
		accessTokenLength:  accessTokenLength,
		refreshTokenLength: refreshTokenLength,
	}
}

func (f *CodeFlowWithOTP) NewStorageUnitOfWork(ctx context.Context) (StorageUnitOfWork, error) {
	return f.storageUoWFactory(ctx)
}

func (f *CodeFlowWithOTP) Start(ctx context.Context, otpReceiverID string) (authCode string, err error) {
	pkcePair, authOperationID, err := f.authProvider.SendOTP(ctx, otpReceiverID)
	if err != nil {
		return "", fmt.Errorf("failed to send otp, err: %w", err)
	}

	return f.codeStorage.Create(ctx, pkcePair, otpReceiverID, authOperationID)
}

func (f *CodeFlowWithOTP) CreateSession(ctx context.Context, authCode, otp string) (*model.Session, error) {
	uow, err := f.NewStorageUnitOfWork(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create new storage unit of work, err: %w", err)
	}

	accessToken := strrand.RandSeqStr(f.accessTokenLength)

	codePayload, err := uow.OAuthCodeStorage().Get(ctx, authCode)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to delete auth code, err: %w", err), uow.Rollback(ctx))
	}

	smartHomeTokens, err := f.authProvider.GetTokensByOTP(
		ctx, codePayload.SmartHomeAuthOperationID, codePayload.SmartHomePKCEPair, otp,
	)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to exchange otp, err: %w", err), uow.Rollback(ctx))
	}

	session, err := uow.SessionStorage().Upsert(
		ctx, "", model.NewSession(
			codePayload.SmartHomeOTPReceiverID, accessToken, strrand.RandSeqStr(f.refreshTokenLength),
			smartHomeTokens.AccessToken, smartHomeTokens.AccessTokenTTL, smartHomeTokens.RefreshToken,
		),
	)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to create new session, err: %w", err), uow.Rollback(ctx))
	}

	if err = uow.OAuthCodeStorage().SetAccessToken(ctx, authCode, accessToken); err != nil {
		return nil, errors.Join(fmt.Errorf("failed to set access token, err: %w", err), uow.Rollback(ctx))
	}

	if err = uow.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit new session, err: %w", err)
	}

	return session, nil
}

func (f *CodeFlowWithOTP) ExchangeAuthCode(ctx context.Context, authCode string) (*model.Session, error) {
	return f.codeStorage.DeleteCodeAndGetSession(ctx, authCode)
}

func (f *CodeFlowWithOTP) GetSessionByAccessToken(ctx context.Context, accessToken string) (*model.Session, error) {
	return f.sessionStorage.GetByAccessToken(ctx, accessToken)
}

func (f *CodeFlowWithOTP) RefreshSession(ctx context.Context, refreshToken string) (*model.Session, error) {
	// TODO: optimize (2 queries -> 1 query)

	session, err := f.sessionStorage.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	tokens, err := f.authProvider.RefreshTokens(ctx, session.SmartHomeRefreshToken)
	if err != nil {
		return nil, err
	}

	refreshed := session.Refresh(
		strrand.RandSeqStr(f.accessTokenLength), strrand.RandSeqStr(f.refreshTokenLength), tokens.AccessToken,
		tokens.AccessTokenTTL, tokens.RefreshToken,
	)

	return f.sessionStorage.Upsert(ctx, session.AccessToken, refreshed)
}
