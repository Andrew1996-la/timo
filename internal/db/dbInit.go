package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DBInit(ctx context.Context) *pgxpool.Pool {
	connString := os.Getenv("CONN_STRING")
	if connString == "" {
		fmt.Println("CONN_STRING is not set")
	}

	// подключение к БД
	pool, err := Connect(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}

	return pool
}
