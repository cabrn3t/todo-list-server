package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/cabrn3t/todo-list-server/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var pool *pgxpool.Pool

func Connect(ctx context.Context, cfg config.Config, connectionAttempts int) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	var err error

	for range connectionAttempts {
		time.Sleep(time.Second * 5)
		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			slog.Error("unable to connect to database...", slog.Any("err", err))
			continue
		}

		if err := pool.Ping(ctx); err != nil {
			slog.Error("error while pinging database", slog.Any("err", err))
			continue
		}

		slog.Info("database successfully connected")
		slog.Info("migrating")

		if err := goose.SetDialect("postgres"); err != nil {
			slog.Error("error while setting dialect", slog.Any("err", err))
			return nil, err
		}

		conn := stdlib.OpenDBFromPool(pool)
		migrationsPath := "./migrations"
		if err := goose.Up(conn, migrationsPath); err != nil {
			slog.Error("error while migrating database", slog.Any("err", err))
			return nil, err
		}

		slog.Info("database migrated successfully")

		return pool, nil
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", connectionAttempts, err)

}
