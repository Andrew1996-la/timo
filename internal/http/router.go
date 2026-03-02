package http

import (
	"net/http"
	"strings"

	"github.com/Andrew1996-la/timo/internal/http/handler"
	"github.com/Andrew1996-la/timo/internal/service"
)

// NewRouter создаёт маршруты для работы с задачами
func NewRouter(taskService *service.TaskService) http.Handler {
	mux := http.NewServeMux()

	// Создаём инстанс TaskHandler, который будет работать с нашим сервисом
	taskHandler := handler.NewTaskHandler(taskService)

	// GET /tasks, POST /tasks
	mux.HandleFunc("/tasks", taskHandler.Tasks)

	// GET /tasks/{id}, DELETE /tasks/{id}, PATCH /tasks/{id}/time
	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch && strings.HasSuffix(r.URL.Path, "/time") {
			taskHandler.AddTime(w, r)
			return
		}
		taskHandler.TaskByID(w, r)
	})

	return mux
}
