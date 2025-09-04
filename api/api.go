package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go-task-manager/internal/task"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// APIHandler holds the API handlers and a Storage instance.
type APIHandler struct {
	Storage task.Storage
}

// Global validator instance
var validate *validator.Validate

func init() {
	// Initialize a new validator instance on package load.
	validate = validator.New()
}

// writeError sends a structured JSON error response.
// It's a helper function to ensure all API error responses are consistent.
func writeError(w http.ResponseWriter, status int, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]interface{}{
		"status":  status,
		"message": message,
		"details": details,
	}
	json.NewEncoder(w).Encode(response)
}

// AddTaskHandler handles requests to add a new task.
func (h *APIHandler) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t task.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		writeError(w, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	err = validate.Struct(t)
	if err != nil {
		log.Error().Err(err).Msg("Validation failed for task data")
		validationErrors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors[e.Field()] = e.Error()
		}
		writeError(w, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	createdTask, err := h.Storage.AddTask(t)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add task to database")
		writeError(w, http.StatusInternalServerError, "Failed to create task", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

// GetTasksHandler handles requests to get a list of tasks.
func (h *APIHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Storage.GetTasks()
	if err != nil {
		// Use writeError for consistent server-side error responses
		log.Error().Err(err).Msg("Failed to retrieve tasks from database")
		writeError(w, http.StatusInternalServerError, "Failed to retrieve tasks", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GetTaskByIDHandler handles requests to retrieve a single task by ID.
func (h *APIHandler) GetTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// Use writeError for invalid client input
		log.Error().Err(err).Msg("Invalid task ID format")
		writeError(w, http.StatusBadRequest, "Invalid task ID", nil)
		return
	}

	task, err := h.Storage.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Use writeError for a consistent "not found" response
			log.Warn().Msgf("Task with ID %d not found", id)
			writeError(w, http.StatusNotFound, "Task not found", nil)
			return
		}
		// Use writeError for server-side database errors
		log.Error().Err(err).Msg("Failed to retrieve task from database")
		writeError(w, http.StatusInternalServerError, "Failed to retrieve task", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *APIHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Err(err).Msg("Invalid task ID format for update")
		writeError(w, http.StatusBadRequest, "Invalid task ID", nil)
		return
	}

	var updatedTask task.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		log.Error().Err(err).Msg("Failed to decode update request body")
		writeError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Add validation for the updated task data
	err = validate.Struct(updatedTask)
	if err != nil {
		log.Error().Err(err).Msg("Validation failed for updated task data")
		// Convert validation errors to a user-friendly format
		validationErrors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors[e.Field()] = e.Error()
		}
		writeError(w, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	result, err := h.Storage.UpdateTask(id, updatedTask)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msgf("Task with ID %d not found for update", id)
			writeError(w, http.StatusNotFound, "Task not found", nil)
			return
		}
		log.Error().Err(err).Msg("Failed to update task in database")
		writeError(w, http.StatusInternalServerError, "Failed to update task", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// DeleteTaskHandler handles requests to delete a task.
func (h *APIHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// Use writeError for invalid client input
		log.Error().Err(err).Msg("Invalid task ID format for deletion")
		writeError(w, http.StatusBadRequest, "Invalid task ID", nil)
		return
	}

	if err := h.Storage.DeleteTask(id); err != nil {
		// Use writeError for server-side database errors
		log.Error().Err(err).Msg("Failed to delete task from database")
		writeError(w, http.StatusInternalServerError, "Failed to delete task", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchTasksHandler handles requests to search for tasks.
func (h *APIHandler) SearchTasksHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		// Use writeError for invalid client input (missing parameter)
		log.Warn().Msg("Search query 'q' is missing")
		writeError(w, http.StatusBadRequest, "Query parameter 'q' is required", nil)
		return
	}

	tasks, err := h.Storage.SearchTasks(query)
	if err != nil {
		// Use writeError for server-side database errors
		log.Error().Err(err).Msg("Failed to search tasks in database")
		writeError(w, http.StatusInternalServerError, "Failed to search for tasks", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
