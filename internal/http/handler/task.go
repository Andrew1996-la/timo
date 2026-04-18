package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Andrew1996-la/timo/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
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
	if r.Method != http.MethodPost {
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

	if req.Seconds <= 0 {
		writeError(w, http.StatusBadRequest, "seconds must be positive")
		return
	}

	if err := h.service.AddTime(r.Context(), id, req.Seconds); err != nil {
		writeServiceError(w, err)
		return
	}

	task, err := h.service.GetById(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) getAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
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
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) getByID(w http.ResponseWriter, r *http.Request, id int) {
	task, err := h.service.GetById(r.Context(), id)
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
	idStr := strings.TrimPrefix(path, "/tasks/")
	if idStr == "" || strings.Contains(idStr, "/") {
		return 0, errors.New("invalid path")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func parseTaskIDFromTimePath(path string) (int, error) {
	if !strings.HasSuffix(path, "/time") {
		return 0, errors.New("invalid path")
	}

	idStr := strings.TrimPrefix(path, "/tasks/")
	idStr = strings.TrimSuffix(idStr, "/time")
	if idStr == "" || strings.Contains(idStr, "/") {
		return 0, errors.New("invalid path")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeMethodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func writeServiceError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusBadRequest, err.Error())
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