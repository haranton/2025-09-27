package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"webfilehosting/internal/models"

	"github.com/google/uuid"
)

type FileStorage struct {
	basePath     string
	tasksDir     string
	downloadsDir string
}

func (fs *FileStorage) SaveTask(task *models.Task) error {
	// Создаем json в tasks присваиваем uid
	jsonData, err := json.MarshalIndent(task, "", " ")
	if err != nil {
		return fmt.Errorf("ошибка преобразования json файла: %w", err)
	}

	fileName := "task_" + uuid.NewString() + ".json"
	filePath := filepath.Join(fs.basePath, fs.tasksDir, fileName)

	if err = os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}
	return nil

}
func (fs *FileStorage) GetTask(taskID string) (*models.Task, error) {
	fileName := "task_" + taskID + ".json"
	filePath := filepath.Join(fs.basePath, fs.tasksDir, fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	var task models.Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return &task, nil

}

func (fs *FileStorage) GetAllTasks() ([]*models.Task, error) {
	dirPath := filepath.Join(fs.basePath, fs.tasksDir)
	files, err := os.ReadDir(dirPath)

	if err != nil {
		return nil, fmt.Errorf("ошибка чтения директории: %w", err)
	}

	var tasks []*models.Task

	for _, file := range files {

		data, err := os.ReadFile(filepath.Join(dirPath, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения файла %s: %w", file.Name(), err)
		}

		var task models.Task
		if err := json.Unmarshal(data, &task); err != nil {
			return nil, fmt.Errorf("ошибка парсинга файла %s: %w", file.Name(), err)
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}
