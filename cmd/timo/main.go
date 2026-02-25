package main

import (
	"context"
	"log"

	"github.com/Andrew1996-la/timo/internal/app"
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

	p := tea.NewProgram(app.New(ctx, taskService))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
