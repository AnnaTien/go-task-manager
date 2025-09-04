package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-task-manager/api"
	"go-task-manager/internal/task"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB initializes an in-memory SQLite database for testing.
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the test database")
	}
	db.AutoMigrate(&task.Task{})
	return db
}

// TestSearchTasksHandler tests the search API endpoint.
func TestSearchTasksHandler(t *testing.T) {
	// Setup an in-memory test database and storage.
	db := setupTestDB()
	storage := task.NewGormStorage(db)
	handler := api.APIHandler{Storage: storage}

	// Add some dummy tasks for testing.
	storage.AddTask(task.Task{Name: "Hoc ve Docker"})
	storage.AddTask(task.Task{Name: "Hoc ve Kubernetes"})
	storage.AddTask(task.Task{Name: "Lap trinh Go co ban"})

	// Create an HTTP request to the search endpoint with a query parameter.
	req, err := http.NewRequest("GET", "/tasks/search?q=Hoc", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a recorder to capture the response.
	rr := httptest.NewRecorder()

	// Create a new router for testing.
	router := http.NewServeMux()
	router.HandleFunc("/tasks/search", handler.SearchTasksHandler)
	router.ServeHTTP(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Decode the response body.
	var tasks []task.Task
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Fatal(err)
	}

	// Verify that the correct number of tasks were returned.
	if len(tasks) != 2 {
		t.Errorf("handler returned unexpected number of tasks: got %v want %v",
			len(tasks), 2)
	}

	// Verify the content of the returned tasks.
	if tasks[0].Name != "Hoc ve Docker" && tasks[1].Name != "Hoc ve Kubernetes" {
		t.Error("handler returned unexpected tasks")
	}
}
