package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/shuryak/sberdevices/internal/model"
	"github.com/shuryak/sberdevices/internal/pkg/pkce"
	"github.com/shuryak/sberdevices/internal/pkg/postgres"
	"github.com/shuryak/sberdevices/internal/pkg/strrand"
)

type OAuthCodeStorage struct {
	codeLength uint

	exec executor
}

func NewOAuthCodeStorage(pg *postgres.Postgres, codeLength uint) *OAuthCodeStorage {
	return &OAuthCodeStorage{
		codeLength: codeLength,
		exec:       pg.Pool,
	}
}

func (storage *OAuthCodeStorage) Create(
	ctx context.Context, smartHomePKCEPair *pkce.Pair, smartHomeOTPReceiverID, smartHomeAuthOperationID string,
) (string, error) {
	authCode := strrand.RandSeqStr(storage.codeLength)

	if _, err := storage.exec.Exec(
		ctx, oauthCodeCreateSQL,
		authCode, smartHomePKCEPair.CodeVerifier, smartHomePKCEPair.CodeChallenge, smartHomeOTPReceiverID,
		smartHomeAuthOperationID,
	); err != nil {
		return "", err
	}

	return authCode, nil
}

func (storage *OAuthCodeStorage) SetAccessToken(ctx context.Context, authCode, accessToken string) error {
	_, err := storage.exec.Exec(ctx, oauthCodeSetAccessTokenSQL, accessToken, authCode)
	return err
}

func (storage *OAuthCodeStorage) Get(
	ctx context.Context, authCode string,
) (*model.AuthCodePayload, error) {
	dbEntity := new(dbSmartHomeOAuthCode)

	// TODO: error not found
	if err := dbEntity.scan(storage.exec.QueryRow(ctx, oauthCodeGetSQL, authCode)); err != nil {
		return nil, err
	}

	return dbEntity.toModel(), nil
}

// DeleteCodeAndGetSession TODO: нарушается модульность. OAuthCodeStorage начинает знать о SessionStorage
func (storage *OAuthCodeStorage) DeleteCodeAndGetSession(ctx context.Context, authCode string) (*model.Session, error) {
	dbEntity := new(dbSession)

	// TODO: error not found
	if err := dbEntity.scan(storage.exec.QueryRow(ctx, oauthCodeDeleteRowAndGetSessionSQL, authCode)); err != nil {
		return nil, err
	}

	return dbEntity.toModel(), nil
}

type dbSmartHomeOAuthCode struct {
	AuthCode                   string
	AccessToken                sql.NullString
	SmartHomePKCECodeVerifier  string
	SmartHomePKCECodeChallenge string
	SmartHomeOTPReceiverID     string
	SmartHomeAuthOperationID   string
}

func (dbEntity *dbSmartHomeOAuthCode) scan(row pgx.Row) error {
	return row.Scan(
		&dbEntity.AuthCode, &dbEntity.AccessToken, &dbEntity.SmartHomePKCECodeVerifier,
		&dbEntity.SmartHomePKCECodeChallenge, &dbEntity.SmartHomeOTPReceiverID, &dbEntity.SmartHomeAuthOperationID,
	)
}

func (dbEntity *dbSmartHomeOAuthCode) toModel() *model.AuthCodePayload {
	accessToken := ""
	if dbEntity.AccessToken.Valid {
		accessToken = dbEntity.AccessToken.String
	}

	return &model.AuthCodePayload{
		SmartHomePKCEPair: &pkce.Pair{
			CodeVerifier:  dbEntity.SmartHomePKCECodeVerifier,
			CodeChallenge: dbEntity.SmartHomePKCECodeChallenge,
		},
		AccessToken:              accessToken,
		SmartHomeOTPReceiverID:   dbEntity.SmartHomeOTPReceiverID,
		SmartHomeAuthOperationID: dbEntity.SmartHomeAuthOperationID,
	}
}
