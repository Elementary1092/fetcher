package postgre

import (
	"context"
	"fmt"
	"github.com/elem1092/fetcher/internal/config"
	"github.com/elem1092/fetcher/pkg/logging"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Close()
}

func NewClient(ctx context.Context, cfg config.DBConfig) (pool *pgxpool.Pool, err error) {
	logger := logging.GetLogger()
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password, cfg.Address, cfg.Port, cfg.DBName)
	logger.Infof("Connecting to the database. URL: %s", connStr)

	for i := 0; i < cfg.MaxAttempts; i++ {
		tryToConnect(func() error {
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			pool, err = pgxpool.Connect(ctx, connStr)
			if err != nil {
				logger.Warnf("Failed to connect due to: %v", err)
				return err
			}

			return nil
		})
	}

	return
}

func tryToConnect(fn func() error) error {
	if err := fn(); err != nil {
		return err
	}

	return nil
}
