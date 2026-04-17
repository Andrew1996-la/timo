package main

import (
	"context"
	"flag"
	"log"

	"github.com/Andrew1996-la/timo/internal/app"
	httptransport "github.com/Andrew1996-la/timo/internal/http"
	"github.com/Andrew1996-la/timo/internal/repository"
	"github.com/Andrew1996-la/timo/internal/service"
	"github.com/Andrew1996-la/timo/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

type container struct {
	taskService *service.TaskService
	cleanup     func()
}

func main() {
	ctx := context.Background()

	c, err := newContainer()
	if err != nil {
		log.Fatal(err)
	}
	defer c.cleanup()

	mode := parseMode()

	switch mode {
	case "http":
		if err := runHTTP(c.taskService); err != nil {
			log.Fatal(err)
		}
	case "cli":
		if err := runCLI(ctx, c.taskService); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown mode: %s", mode)
	}
}

func newContainer() (*container, error) {
	db, err := storage.NewSQLite()
	if err != nil {
		return nil, err
	}

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)

	return &container{
		taskService: taskService,
		cleanup: func() {
			if err := db.Close(); err != nil {
				log.Printf("failed to close db: %v", err)
			}
		},
	}, nil
}

func parseMode() string {
	httpMode := flag.Bool("http", false, "run HTTP server")
	cliMode := flag.Bool("cli", false, "run CLI app")
	flag.Parse()

	switch {
	case *httpMode:
		return "http"
	case *cliMode:
		return "cli"
	default:
		return "cli"
	}
}

func runHTTP(taskService *service.TaskService) error {
	router := httptransport.NewRouter(taskService)
	server := httptransport.New(":8080", router)

	log.Println("HTTP server started on :8080")
	return server.Start()
}

func runCLI(ctx context.Context, taskService *service.TaskService) error {
	p := tea.NewProgram(app.New(ctx, taskService))
	_, err := p.Run()
	return err
}