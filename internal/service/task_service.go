package service

import (
	"webfilehosting/internal/models"
	"webfilehosting/internal/storage"

	"github.com/google/uuid"
)

type TaskService struct {
	storage *storage.FileStorage
}

// Добавляем конструктор
func NewTaskService(storage *storage.FileStorage) *TaskService {
	return &TaskService{
		storage: storage,
	}
}

func (s *TaskService) CreateTask(urls []string) (*models.Task, error) {
	fileInfos := make([]models.FileInfo, len(urls))
	for i, u := range urls {
		fileInfos[i] = models.FileInfo{
			URL:    u,
			Status: models.FileStatusPending,
		}
	}
	task := models.Task{
		Id:     uuid.New(),
		Status: models.TaskStatusPending,
		Urls:   fileInfos,
	}
	if err := s.storage.SaveTask(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) GetTask(taskID string) (*models.Task, error) {
	return s.storage.GetTask(taskID)
}

func (s *TaskService) GetAllTasks() ([]*models.Task, error) {
	return s.storage.GetAllTasks()
}

func (s *TaskService) GetPendingTasks() ([]*models.Task, error) {
	all, err := s.storage.GetAllTasks()
	if err != nil {
		return nil, err
	}
	var pending []*models.Task
	for _, t := range all {
		if t.Status == models.TaskStatusPending {
			pending = append(pending, t)
		}
	}
	return pending, nil
}

func (s *TaskService) UpdateTaskStatus(taskID string, status models.TaskStatus) error {
	return s.storage.UpdateTaskStatus(taskID, status)
}

func (s *TaskService) UpdateFileInfo(taskID string, fileInfo *models.FileInfo) error {
	task, err := s.GetTask(taskID)
	if err != nil {
		return err
	}
	for i := range task.Urls {
		if task.Urls[i].URL == fileInfo.URL {
			task.Urls[i] = *fileInfo
			break
		}
	}
	return s.storage.SaveTask(task)
}

func (s *TaskService) SaveTask(task *models.Task) error {
	return s.storage.SaveTask(task)
}

// Добавляем метод для захвата задачи (чтобы избежать гонок)
func (s *TaskService) AcquireTask(taskID string) (bool, error) {
	task, err := s.GetTask(taskID)
	if err != nil {
		return false, err
	}

	if task.Status != models.TaskStatusPending {
		return false, nil // Задача уже кем-то взята
	}

	// Меняем статус на "в процессе"
	err = s.UpdateTaskStatus(taskID, models.TaskStatusInProgress)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Добавляем метод для сброса задачи
func (s *TaskService) ResetTask(taskID string) error {
	task, err := s.GetTask(taskID)
	if err != nil {
		return err
	}

	// Сбрасываем статус задачи
	task.Status = models.TaskStatusPending

	// Сбрасываем статусы файлов
	for i := range task.Urls {
		if task.Urls[i].Status == models.FileStatusDownloading {
			task.Urls[i].Status = models.FileStatusPending
			task.Urls[i].Error = ""
		}
	}

	return s.storage.SaveTask(task)
}
