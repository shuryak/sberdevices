package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/shuryak/sberdevices/internal/oauth"
	"github.com/shuryak/sberdevices/internal/pkg/postgres"
)

type UnitOfWork struct {
	sessionStorage   *SessionStorage
	oAuthCodeStorage *OAuthCodeStorage

	tx pgx.Tx
}

func NewUnitOfWork(
	ctx context.Context, pg *postgres.Postgres, sessionStorage *SessionStorage, oAuthCodeStorage *OAuthCodeStorage,
) (*UnitOfWork, error) {
	tx, err := pg.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	txSession := *sessionStorage
	txSession.exec = tx

	txOAuthCodeStorage := *oAuthCodeStorage
	txOAuthCodeStorage.exec = tx

	return &UnitOfWork{
		tx:               tx,
		sessionStorage:   &txSession,
		oAuthCodeStorage: &txOAuthCodeStorage,
	}, nil
}

func NewUnitOfWorkFactory(
	pg *postgres.Postgres, sessionStorage *SessionStorage, oAuthCodeStorage *OAuthCodeStorage,
) oauth.StorageUnitOfWorkFactory {
	return func(ctx context.Context) (oauth.StorageUnitOfWork, error) {
		return NewUnitOfWork(ctx, pg, sessionStorage, oAuthCodeStorage)
	}
}

func (uow *UnitOfWork) SessionStorage() oauth.SessionStorage {
	return uow.sessionStorage
}

func (uow *UnitOfWork) OAuthCodeStorage() oauth.CodeStorage {
	return uow.oAuthCodeStorage
}

func (uow *UnitOfWork) Commit(ctx context.Context) error {
	return uow.tx.Commit(ctx)
}

func (uow *UnitOfWork) Rollback(ctx context.Context) error {
	return uow.tx.Rollback(ctx)
}
