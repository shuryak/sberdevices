package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shuryak/sberdevices/internal/model"
	"github.com/shuryak/sberdevices/internal/pkg/postgres"
)

type SessionStorage struct {
	exec executor
}

func NewSessionStorage(pg *postgres.Postgres) *SessionStorage {
	return &SessionStorage{
		exec: pg.Pool,
	}
}

func (storage *SessionStorage) Upsert(
	ctx context.Context, oldAccessToken string, session *model.Session,
) (*model.Session, error) {
	storedSession := new(dbSession)

	if len(oldAccessToken) == 0 {
		oldAccessToken = session.AccessToken
	}

	if err := storedSession.scan(storage.exec.QueryRow(ctx, sessionUpsertSQL,
		oldAccessToken, session.AccessToken, session.RefreshToken, session.OTPReceiverID, session.SmartHomeAccessToken,
		session.SmartHomeAccessTokenTTL, session.SmartHomeRefreshToken, session.CreatedAt, session.UpdatedAt,
	)); err != nil {
		return nil, err
	}

	return storedSession.toModel(), nil
}

func (storage *SessionStorage) GetByAccessToken(ctx context.Context, accessToken string) (*model.Session, error) {
	session := new(dbSession)

	if err := session.scan(storage.exec.QueryRow(ctx, sessionGetByAccessTokenSQL, accessToken)); err != nil {
		return nil, err
	}

	return session.toModel(), nil
}

func (storage *SessionStorage) GetByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error) {
	session := new(dbSession)

	if err := session.scan(storage.exec.QueryRow(ctx, sessionGetByRefreshTokenSQL, refreshToken)); err != nil {
		return nil, err
	}

	return session.toModel(), nil
}

type dbSession struct {
	AccessToken  string
	RefreshToken string

	OTPReceiverID string

	SmartHomeAccessToken    string
	SmartHomeAccessTokenTTL time.Duration

	SmartHomeRefreshToken string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (dbSession *dbSession) scan(row pgx.Row) error {
	return row.Scan(
		&dbSession.AccessToken, &dbSession.RefreshToken, &dbSession.OTPReceiverID, &dbSession.SmartHomeAccessToken,
		&dbSession.SmartHomeAccessTokenTTL, &dbSession.SmartHomeRefreshToken, &dbSession.CreatedAt,
		&dbSession.UpdatedAt,
	)
}

func (dbSession *dbSession) toModel() *model.Session {
	return &model.Session{
		OTPReceiverID:           dbSession.OTPReceiverID,
		AccessToken:             dbSession.AccessToken,
		RefreshToken:            dbSession.RefreshToken,
		SmartHomeAccessToken:    dbSession.SmartHomeAccessToken,
		SmartHomeAccessTokenTTL: dbSession.SmartHomeAccessTokenTTL,
		SmartHomeRefreshToken:   dbSession.SmartHomeRefreshToken,
		UpdatedAt:               dbSession.UpdatedAt,
		CreatedAt:               dbSession.CreatedAt,
	}
}
