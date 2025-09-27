package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config представляет конфигурацию приложения
type Config struct {
	DirectoryDownload string
	StorageFileJSON   string
	Host              string
	Port              int
}

var AppConfig *Config

// InitConfig инициализирует конфигурацию из .env файла
func InitConfig() error {
	config := &Config{
		DirectoryDownload: "uploads",
		StorageFileJSON:   "tasks",
		Host:              "localhost",
		Port:              8080,
	}

	// Загружаем переменные из .env файла
	if err := loadFromEnvFile(config); err != nil {
		fmt.Printf("Предупреждение: не удалось загрузить .env файл: %v\n", err)
	}

	// Валидируем конфигурацию
	if err := validateConfig(config); err != nil {
		return fmt.Errorf("ошибка валидации конфигурации: %v", err)
	}

	// Создаем необходимые директории
	if err := createDirectories(config); err != nil {
		return fmt.Errorf("ошибка создания директорий: %v", err)
	}

	// Сохраняем конфигурацию в глобальную переменную
	AppConfig = config

	return nil
}

// loadFromEnvFile загружает переменные из .env файла
func loadFromEnvFile(config *Config) error {
	envFile := ".env"
	
	file, err := os.Open(envFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Пропускаем пустые строки и комментарии
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Разбираем строку KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Убираем кавычки если есть
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || 
			(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}

		// Устанавливаем значения в зависимости от ключа
		switch key {
		case "DIRECTORY_DOWNLOAD":
			config.DirectoryDownload = value
		case "STORAGEFILE_JSON":
			config.StorageFileJSON = value
		case "HOST":
			config.Host = value
		case "PORT":
			if port, err := strconv.Atoi(value); err == nil {
				config.Port = port
			}
		}
	}

	return scanner.Err()
}

// validateConfig проверяет корректность конфигурации
func validateConfig(config *Config) error {
	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("неверный порт: %d", config.Port)
	}

	if config.DirectoryDownload == "" {
		return fmt.Errorf("директория загрузки не может быть пустой")
	}

	if config.StorageFileJSON == "" {
		return fmt.Errorf("имя JSON файла не может быть пустым")
	}

	if config.Host == "" {
		return fmt.Errorf("хост не может быть пустым")
	}

	return nil
}

// createDirectories создает необходимые директории
func createDirectories(config *Config) error {
	// Создаем директорию для загрузок
	if err := os.MkdirAll(config.DirectoryDownload, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию загрузок: %v", err)
	}

	return nil
}

// GetServerAddress возвращает адрес сервера в формате host:port
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetUploadPath возвращает полный путь для загрузки файла
func (c *Config) GetUploadPath(filename string) string {
	return filepath.Join(c.DirectoryDownload, filename)
}

// GetStorageFilePath возвращает полный путь к JSON файлу
func (c *Config) GetStorageFilePath() string {
	return c.StorageFileJSON + ".json"
}