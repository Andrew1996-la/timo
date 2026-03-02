package main

import (
	"context"
	"log"

	"github.com/Andrew1996-la/timo/internal/app"
	"github.com/Andrew1996-la/timo/internal/http"
	"github.com/Andrew1996-la/timo/internal/repository"
	"github.com/Andrew1996-la/timo/internal/service"
	"github.com/Andrew1996-la/timo/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ctx := context.Background()

	db, err := storage.NewSQLite()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// инициализация репозитория
	taskRepo := repository.NewTaskRepository(db)
	// инициализация сервиса
	taskService := service.NewTaskService(taskRepo)

	router := http.NewRouter(taskService)
	server := http.New(":8080", router)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Println("HTTP server started on :8080")
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// Запуск терминального приложения
	p := tea.NewProgram(app.New(ctx, taskService))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
