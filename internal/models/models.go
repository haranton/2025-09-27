package models

import "github.com/google/uuid"

type TaskStatus string
type FileStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusError      TaskStatus = "error"

	FileStatusPending     FileStatus = "pending"
	FileStatusDownloading FileStatus = "downloading"
	FileStatusCompleted   FileStatus = "completed"
	FileStatusError       FileStatus = "error"
)

type Task struct {
	Id     uuid.UUID
	Status TaskStatus
	Urls   []string
}

type File struct {
	URL      string     `json:"url"` // Естественный идентификатор
	Status   FileStatus `json:"status"`
	FilePath string     `json:"file_path,omitempty"`
	Error    string     `json:"error,omitempty"`
}
