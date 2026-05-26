package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Andrew1996-la/timo/internal/app"
	httptransport "github.com/Andrew1996-la/timo/internal/http"
	"github.com/Andrew1996-la/timo/internal/repository"
	"github.com/Andrew1996-la/timo/internal/service"
	"github.com/Andrew1996-la/timo/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	modeHTTP = "http"
	modeCLI  = "cli"

	shutdownTimeout = 5 * time.Second
)

type config struct {
	mode string
	addr string
}

type container struct {
	taskService *service.TaskService
	cleanup     func()
}

func main() {
	cfg := parseConfig()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	c, err := newContainer()
	if err != nil {
		log.Fatal(err)
	}
	defer c.cleanup()

	if err := run(ctx, cfg, c); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, cfg config, c *container) error {
	switch cfg.mode {
	case modeHTTP:
		return runHTTP(ctx, cfg.addr, c.taskService)
	case modeCLI:
		return runCLI(ctx, c.taskService)
	default:
		return errors.New("unknown mode: " + cfg.mode)
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

func parseConfig() config {
	httpMode := flag.Bool("http", false, "run HTTP server")
	cliMode := flag.Bool("cli", false, "run CLI app")
	addr := flag.String("addr", ":8080", "HTTP server address")

	flag.Parse()

	mode := modeCLI
	if *httpMode {
		mode = modeHTTP
	}

	if *cliMode {
		mode = modeCLI
	}

	return config{
		mode: mode,
		addr: *addr,
	}
}

func runHTTP(
	ctx context.Context,
	addr string,
	taskService *service.TaskService,
) error {
	router := httptransport.NewRouter(taskService)
	server := httptransport.New(addr, router)

	errCh := make(chan error, 1)

	go func() {
		log.Printf("HTTP server started on %s", addr)

		if err := server.Start(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}

		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			shutdownTimeout,
		)
		defer cancel()

		log.Println("shutting down HTTP server")

		return server.Shutdown(shutdownCtx)

	case err := <-errCh:
		return err
	}
}

func runCLI(ctx context.Context, taskService *service.TaskService) error {
	program := tea.NewProgram(app.New(ctx, taskService))

	_, err := program.Run()

	return err
}