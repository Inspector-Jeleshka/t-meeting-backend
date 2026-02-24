package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, dsn string) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse postgres dsn: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	// общий лимит на установление связи (чтобы не висеть вечно)
	connectCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	const (
		pingTimeout = 2 * time.Second
		baseDelay   = 200 * time.Millisecond
		maxDelay    = 2 * time.Second
	)

	var lastErr error
	for attempt := 1; ; attempt++ {
		pingCtx, cancelPing := context.WithTimeout(connectCtx, pingTimeout)
		err = pool.Ping(pingCtx)
		cancelPing()

		if err == nil {
			// успех)
			return &DB{pool: pool}, nil
		}

		lastErr = err

		// если общий таймаут вышел сдаемся
		if connectCtx.Err() != nil {
			pool.Close()
			return nil, fmt.Errorf("ping postgres (attempts=%d): %w", attempt, lastErr)
		}

		delay := baseDelay * time.Duration(1<<(attempt-1))
		if delay > maxDelay {
			delay = maxDelay
		}
		time.Sleep(delay)
	}
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *DB) Close() {
	db.pool.Close()
}
