package client

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, dsn string) *pgxpool.Pool {
	conn, err := pgxpool.New(ctx, dsn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil
	}

	var greeting string
	err = conn.QueryRow(ctx, "select 'Connected to database.\n'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return nil
	}

	fmt.Println(greeting)

	return conn
}
