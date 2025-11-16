package postgres

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil { return nil, err }

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.HealthCheckPeriod = 30 * time.Second
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil { return nil, err }

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, err
	}
	log.Println("postgres: connected")
	return pool, nil
}

func openStdDB(dbURL string) *sql.DB {
	db, err := sql.Open("pgx", dbURL)
	if err != nil { panic(err) }
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	return db
}

func MustPool(ctx context.Context) *pgxpool.Pool {
	url := os.Getenv("DATABASE_URL")
	if url == "" { panic("DATABASE_URL empty") }
	p, err := NewPool(ctx, url)
	if err != nil { panic(err) }
	return p
}