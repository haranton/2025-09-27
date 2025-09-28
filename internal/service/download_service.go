package service

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type DownloadService struct {
	httpClient *http.Client
}

func NewDownloadService() *DownloadService {
	return &DownloadService{
		httpClient: &http.Client{
			Timeout: 30 * time.Minute,
		},
	}
}

func (s *DownloadService) DownloadFile(fileURL, savePath string) error {
	log.Printf("Downloading: %s -> %s\n", fileURL, savePath)

	if err := s.validateURL(fileURL); err != nil {
		return err
	}

	if err := s.ensureDirectory(savePath); err != nil {
		return err
	}

	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	log.Printf("Download completed: %s\n", fileURL)
	return nil
}

func (s *DownloadService) validateURL(fileURL string) error {
	parsed, err := url.Parse(fileURL)
	if err != nil {
		return err
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return err
	}

	if parsed.Host == "" {
		return err
	}

	return nil
}

func (s *DownloadService) ensureDirectory(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}

func (s *DownloadService) CreateTaskDirectory(taskDir string) error {
	log.Printf("Creating task directory: %s\n", taskDir)
	return s.ensureDirectory(filepath.Join(taskDir, "dummy.txt"))
}

func (s *DownloadService) GenerateFileName(fileURL string) string {
	parsed, err := url.Parse(fileURL)
	if err != nil {
		return "file"
	}

	path := parsed.Path
	if path == "" || path == "/" {
		return "file"
	}

	baseName := filepath.Base(path)
	if baseName == "" || baseName == "." || baseName == "/" {
		return "file"
	}

	return baseName
}
