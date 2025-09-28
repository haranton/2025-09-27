package models

import (
	"time"

	"github.com/google/uuid"
)

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

type FileInfo struct {
	URL         string     `json:"url"`
	Status      FileStatus `json:"status"`
	FilePath    string     `json:"file_path,omitempty"`
	Error       string     `json:"error,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type Task struct {
	Id     uuid.UUID  `json:"id"`
	Status TaskStatus `json:"status"`
	Urls   []FileInfo `json:"urls"`
}
