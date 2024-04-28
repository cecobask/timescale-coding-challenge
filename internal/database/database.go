package database

import (
	"context"
	"fmt"
	"os"
	"time"

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

const cpuUsageQuery = `
	SELECT
		time_bucket('1 minutes', ts) AS bucket,
		first(usage, ts),
		last(usage, ts) 
	FROM
		cpu_usage 
	WHERE
		host = $1 AND ts BETWEEN $2 AND $3 
	GROUP BY
		bucket 
	ORDER BY
		bucket ASC;
`

func (db *Database) BenchmarkQuery(ctx context.Context, hostname, startTime, endTime string) (time.Duration, error) {
	start := time.Now()
	if _, err := db.pool.Exec(ctx, cpuUsageQuery, hostname, startTime, endTime); err != nil {
		return 0, fmt.Errorf("unable to execute query: %w", err)
	}
	return time.Since(start), nil
}

func (db *Database) Close() {
	db.pool.Close()
}
