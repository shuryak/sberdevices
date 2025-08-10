package postgres

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(connString string, opts *Options, log *slog.Logger) (*Postgres, error) {
	pg := new(Postgres)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = opts.maxPoolSize

	pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	if opts.connAttempts < 1 {
		return nil, errors.New("connAttempts must be >= 1")
	}

	for i := range opts.connAttempts {
		if err = pg.Pool.Ping(context.Background()); err != nil {
			log.Error(
				"connection to postgres failed", slog.Int("attempts_left", int(opts.connAttempts-i-1)),
				slog.String("error", err.Error()),
			)
			time.Sleep(opts.connTimeout)
		} else {
			log.Info("connection to postgres established")
			return pg, nil
		}
	}

	return pg, err
}

func (p *Postgres) Close() {
	if p != nil && p.Pool != nil {
		p.Pool.Close()
	}
}
