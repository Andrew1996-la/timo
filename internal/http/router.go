package http

import (
	"net/http"

	"github.com/Andrew1996-la/timo/internal/http/handler"
	"github.com/Andrew1996-la/timo/internal/service"
)

func NewRouter(taskService *service.TaskService) http.Handler {
	mux := http.NewServeMux()

	taskHandler := handler.NewTaskHandler(taskService)

	mux.HandleFunc("/tasks", taskHandler.Tasks)
	mux.HandleFunc("/tasks/", taskHandler.TaskByID)

	return mux
}
