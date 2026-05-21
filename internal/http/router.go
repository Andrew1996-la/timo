package http

import (
	"net/http"
	"strings"

	"github.com/Andrew1996-la/timo/internal/http/handler"
	"github.com/Andrew1996-la/timo/internal/service"
)

func NewRouter(taskService *service.TaskService) http.Handler {
	mux := http.NewServeMux()

	taskHandler := handler.NewTaskHandler(taskService)

	mux.HandleFunc("/tasks", taskHandler.Tasks)
	mux.HandleFunc("/tasks/", taskRoute(taskHandler))

	return mux
}

func taskRoute(taskHandler *handler.TaskHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isTaskTimePath(r.URL.Path) {
			taskHandler.AddTime(w, r)
			return
		}

		taskHandler.TaskByID(w, r)
	}
}

func isTaskTimePath(path string) bool {
	return strings.HasSuffix(path, "/time")
}