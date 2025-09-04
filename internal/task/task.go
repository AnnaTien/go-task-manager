package task

import (
	"gorm.io/gorm"
)

// Task represents a single task in the application.
// It includes GORM and JSON tags for database and API interactions.
type Task struct {
	gorm.Model        // GORM automatically adds ID, CreatedAt, UpdatedAt, and DeletedAt fields.
	Name       string `gorm:"size:255;not null" json:"name" validate:"required,min=3,max=255"`
	Completed  bool   `gorm:"default:false" json:"completed"`
}

// Storage is an interface that defines the methods for
// managing tasks.
type Storage interface {
	GetTasks() ([]Task, error)
	GetTaskByID(id int) (Task, error)
	AddTask(Task) (Task, error)
	DeleteTask(id int) error
	UpdateTask(id int, updatedTask Task) (Task, error)
	SearchTasks(query string) ([]Task, error)
}
