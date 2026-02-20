package main

import (
	"context"
	"log"

	"github.com/Andrew1996-la/timo/internal/app"
	"github.com/Andrew1996-la/timo/internal/db"
	"github.com/Andrew1996-la/timo/internal/repository"
	"github.com/Andrew1996-la/timo/internal/service"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ctx := context.Background()

	pool := db.DBInit(ctx)
	defer pool.Close()
	// инициализация репозитория
	taskRepo := repository.NewTaskRepository(pool)
	//// инициализация сервиса
	taskService := service.NewTaskService(taskRepo)

	p := tea.NewProgram(app.New(ctx, taskService))

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
