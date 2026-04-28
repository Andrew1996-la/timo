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
	mux.HandleFunc("/tasks/", routeTaskByID(taskHandler))

	return mux
}

func routeTaskByID(taskHandler *handler.TaskHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isAddTimeRoute(r) {
			taskHandler.AddTime(w, r)
			return
		}

		taskHandler.TaskByID(w, r)
	}
}

func isAddTimeRoute(r *http.Request) bool {
	return r.Method == http.MethodPatch &&
		strings.HasSuffix(r.URL.Path, "/time")
}