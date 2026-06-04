package http

import (
	"net/http"
	"strings"

	"github.com/Andrew1996-la/timo/internal/http/handler"
	"github.com/Andrew1996-la/timo/internal/service"
)

const (
	tasksPath    = "/tasks"
	taskPath     = "/tasks/"
	taskTimePath = "/time"
)

func NewRouter(taskService *service.TaskService) http.Handler {
	mux := http.NewServeMux()
	taskHandler := handler.NewTaskHandler(taskService)

	mux.HandleFunc(tasksPath, taskHandler.Tasks)
	mux.HandleFunc(taskPath, routeTaskByID(taskHandler))

	return mux
}

func routeTaskByID(taskHandler *handler.TaskHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if hasTaskTimeSuffix(r.URL.Path) {
			taskHandler.AddTime(w, r)
			return
		}

		taskHandler.TaskByID(w, r)
	}
}

func hasTaskTimeSuffix(path string) bool {
	return strings.HasSuffix(path, taskTimePath)
}