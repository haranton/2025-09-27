package service

import (
	"webfilehosting/internal/models"
	"webfilehosting/internal/storage"

	"github.com/google/uuid"
)

type TaskService struct {
	storage *storage.FileStorage
}

func (s *TaskService) CreateTask(urls []string) (*models.Task, error) {

	task := models.Task{
		Id:     uuid.New(),
		Status: models.TaskStatusPending,
		Urls:   urls,
	}

	if err := s.storage.SaveTask(&task); err != nil {
		return nil, err
	}

	return &task, nil

}

func (s *TaskService) GetTask(taskID string) (*models.Task, error) {
	return s.storage.GetTask(taskID)
}
func (s *TaskService) GetPendingTasks() ([]*models.Task, error) {
	return s.storage.GetAllTasks()
}
