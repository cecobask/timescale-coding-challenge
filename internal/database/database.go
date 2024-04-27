package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	envVarNameUsername = "POSTGRES_USER"
	envVarNamePassword = "POSTGRES_PASSWORD"
	envVarNameHost     = "POSTGRES_HOST"
	envVarNameDatabase = "POSTGRES_DB"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context) (*Database, error) {
	username := os.Getenv(envVarNameUsername)
	password := os.Getenv(envVarNamePassword)
	host := os.Getenv(envVarNameHost)
	database := os.Getenv(envVarNameDatabase)
	connString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", username, password, host, database)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create database connection pool: %w", err)
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	return &Database{
		pool: pool,
	}, nil
}

func (d *Database) Close() {
	d.pool.Close()
}
