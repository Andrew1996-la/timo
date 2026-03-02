package main

import (
	"context"
	"log"
	"os"

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

	runHTTP := false
	runCLI := false
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--http":
			runHTTP = true
		case "--cli":
			runCLI = true
		}
	}

	if runHTTP {
		router := http.NewRouter(taskService)
		server := http.New(":8080", router)
		log.Println("HTTP server started on :8080")
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if runCLI {
		p := tea.NewProgram(app.New(ctx, taskService))
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// По умолчанию можно запускать CLI
	p := tea.NewProgram(app.New(ctx, taskService))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
