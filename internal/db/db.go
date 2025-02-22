package db

import (
	"context"
	"electronic-library/config"
	"fmt"
	"os"
	"time"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context, cfg *config.Config) (*Database, error) {
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is not set")
	}

	pgxConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %v", err)
	}

	pgxConfig.ConnConfig.Config.Host = cfg.Database.Host
	pgxConfig.ConnConfig.Config.Port = cfg.Database.Port
	pgxConfig.ConnConfig.Config.User = cfg.Database.User
	pgxConfig.ConnConfig.Config.Password = dbPassword
	pgxConfig.MaxConns = cfg.Database.MaxConns
	pgxConfig.MinConns = cfg.Database.MinConns

	maxConnLifetime, err := time.ParseDuration(cfg.Database.MaxConnLifetime)
	if err != nil {
		return nil, fmt.Errorf("invalid MaxConnLifetime value: %v", err)
	}
	pgxConfig.MaxConnLifetime = maxConnLifetime

	maxConnIdleTime, err := time.ParseDuration(cfg.Database.MaxConnIdleTime)
	if err != nil {
		return nil, fmt.Errorf("invalid MaxConnIdleTime value: %v", err)
	}
	pgxConfig.MaxConnIdleTime = maxConnIdleTime

	healthCheckPeriod, err := time.ParseDuration(cfg.Database.HealthCheckPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid HealthCheckPeriod value: %v", err)
	}
	pgxConfig.HealthCheckPeriod = healthCheckPeriod

	maxConnLifetimeJitter, err := time.ParseDuration(cfg.Database.MaxConnLifetimeJitter)
	if err != nil {
		return nil, fmt.Errorf("invalid MaxConnLifetimeJitter value: %v", err)
	}
	pgxConfig.MaxConnLifetimeJitter = maxConnLifetimeJitter

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	connection, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the database: %v", err)
	}

	return &Database{Pool: connection}, nil
}

func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}
