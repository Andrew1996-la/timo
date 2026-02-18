package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Andrew1996-la/timo/internal/db"
)

func main() {
	ctx := context.Background()
	connString := os.Getenv("CONN_STRING")
	if connString == "" {
		fmt.Println("CONN_STRING is not set")
	}

	pool, err := db.Connect(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	log.Println("Connected to database")
}
