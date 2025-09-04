package task

import (
	"gorm.io/gorm"
)

// GormStorage is a struct that implements the Storage interface
// using the GORM ORM.
type GormStorage struct {
	db *gorm.DB
}

// NewGormStorage creates a new instance of GormStorage.
func NewGormStorage(db *gorm.DB) *GormStorage {
	return &GormStorage{db: db}
}

// GetTasks retrieves all tasks from the database.
func (s *GormStorage) GetTasks() ([]Task, error) {
	var tasks []Task
	result := s.db.Find(&tasks)
	return tasks, result.Error
}

// AddTask adds a new task to the database.
func (s *GormStorage) AddTask(t Task) (Task, error) {
	result := s.db.Create(&t)
	return t, result.Error
}

// GetTaskByID retrieves a single task by its ID from the database.
func (s *GormStorage) GetTaskByID(id int) (Task, error) {
	var t Task
	result := s.db.First(&t, id)
	return t, result.Error
}

// DeleteTask removes a task from the database by its ID.
func (s *GormStorage) DeleteTask(id int) error {
	result := s.db.Delete(&Task{}, id)
	return result.Error
}

// UpdateTask updates an existing task in the database.
func (s *GormStorage) UpdateTask(id int, updatedTask Task) (Task, error) {
	var t Task
	result := s.db.First(&t, id)
	if result.Error != nil {
		return t, result.Error
	}
	t.Name = updatedTask.Name
	t.Completed = updatedTask.Completed
	result = s.db.Save(&t)
	return t, result.Error
}

// SearchTasks finds tasks in the database based on a query string.
func (s *GormStorage) SearchTasks(query string) ([]Task, error) {
	var tasks []Task
	// Use SQL's LIKE command to find a partial match in the task name.
	// The '%' is a wildcard character.
	result := s.db.Where("name LIKE ?", "%"+query+"%").Find(&tasks)
	return tasks, result.Error
}
