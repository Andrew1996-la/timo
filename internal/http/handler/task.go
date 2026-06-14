package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Andrew1996-la/timo/internal/models"
	"github.com/Andrew1996-la/timo/internal/repository"
	"github.com/Andrew1996-la/timo/internal/service"
)

const (
	tasksPathPrefix = "/tasks/"
	timePathSuffix  = "/time"
)

type taskService interface {
	GetAll(ctx context.Context) ([]models.Task, error)
	Create(ctx context.Context, title string) (*models.Task, error)
	GetByID(ctx context.Context, id int) (*models.Task, error)
	Delete(ctx context.Context, id int) error
	AddTime(ctx context.Context, id int, seconds int) error
}

type TaskHandler struct {
	service taskService
}

func NewTaskHandler(service taskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

type createTaskRequest struct {
	Title string `json:"title"`
}

type addTimeRequest struct {
	Seconds int `json:"seconds"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		writeMethodNotAllowed(w)
	}
}

func (h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseTaskID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getByID(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		writeMethodNotAllowed(w)
	}
}

func (h *TaskHandler) AddTime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		writeMethodNotAllowed(w)
		return
	}

	id, err := parseTaskIDFromTimePath(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req addTimeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.service.AddTime(r.Context(), id, req.Seconds); err != nil {
		writeServiceError(w, err)
		return
	}

	task, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) getAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get tasks")
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	task, err := h.service.Create(r.Context(), req.Title)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) getByID(w http.ResponseWriter, r *http.Request, id int) {
	task, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseTaskID(path string) (int, error) {
	return parseIDFromPath(path, tasksPathPrefix, "")
}

func parseTaskIDFromTimePath(path string) (int, error) {
	return parseIDFromPath(path, tasksPathPrefix, timePathSuffix)
}

func parseIDFromPath(path string, prefix string, suffix string) (int, error) {
	value := strings.TrimPrefix(path, prefix)

	if suffix != "" {
		if !strings.HasSuffix(value, suffix) {
			return 0, errors.New("invalid path")
		}

		value = strings.TrimSuffix(value, suffix)
	}

	if value == "" || strings.Contains(value, "/") {
		return 0, errors.New("invalid path")
	}

	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	if id <= 0 {
		return 0, errors.New("id must be positive")
	}

	return id, nil
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(dst)
}

func writeMethodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrTaskNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrInvalidTaskID):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrEmptyTaskTitle):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrInvalidDuration):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{
		Error: message,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(v)
}
