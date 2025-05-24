package db

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPgPool(uri string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	trimmed := strings.Trim(strings.TrimSpace(uri), "\n")
	pool, err := pgxpool.Connect(ctx, trimmed)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
