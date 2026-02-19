package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Andrew1996-la/timo/internal/db"
	"github.com/Andrew1996-la/timo/internal/repository"
	"github.com/Andrew1996-la/timo/internal/service"
)

func main() {
	ctx := context.Background()
	connString := os.Getenv("CONN_STRING")
	if connString == "" {
		fmt.Println("CONN_STRING is not set")
	}

	// подключение к БД
	pool, err := db.Connect(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// инициализация репозитория
	taskRepo := repository.NewTaskRepository(pool)

	// инициализация сервиса
	taskService := service.NewTaskService(taskRepo)

	_, err = taskService.Create(ctx, "Получить оффер 300 тысяч рублей")
	if err != nil {
		log.Fatal(err)
	}

}
