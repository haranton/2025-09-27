package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"webfilehosting/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализируем конфигурацию
	if err := config.InitConfig(); err != nil {
		fmt.Printf("Ошибка инициализации конфигурации: %v\n", err)
		os.Exit(1)
	}

	r := gin.Default()
	
	// Настраиваем CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Web File Hosting Service",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	r.POST("/uploads", CreateDownload)

	// Запускаем сервер с адресом из конфигурации
	serverAddr := config.AppConfig.GetServerAddress()
	fmt.Printf("Сервер запущен на %s\n", serverAddr)
	r.Run(serverAddr)
}

func CreateDownload(ctx *gin.Context) {
	var body struct {
		Url string `json:"url"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON", "details": err.Error()})
		return
	}

	// Проверяем URL
	if body.Url == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "URL не может быть пустым"})
		return
	}

	// Создаем HTTP клиент с ограничением размера
	client := &http.Client{}
	req, err := http.NewRequest("GET", body.Url, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный URL", "details": err.Error()})
		return
	}

	// Добавляем User-Agent для избежания блокировки
	req.Header.Set("User-Agent", "WebFileHosting/1.0")

	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка загрузки файла", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "Ошибка загрузки файла",
			"status_code": resp.StatusCode,
			"status":      resp.Status,
		})
		return
	}

	// Получаем Content-Type для информации
	contentType := resp.Header.Get("Content-Type")

	// Извлекаем имя файла из URL
	filename := filepath.Base(body.Url)
	if filename == "" || filename == "." || filename == "/" {
		filename = "downloaded_file"
	}

	// Создаем путь к файлу используя конфигурацию
	filePath := config.AppConfig.GetUploadPath(filename)

	// Создаем файл
	out, err := os.Create(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания файла", "details": err.Error()})
		return
	}
	defer out.Close()

	// Копируем данные
	bytesWritten, err := io.Copy(out, resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения файла", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Файл успешно загружен",
		"filename":     filename,
		"path":         filePath,
		"size":         bytesWritten,
		"content_type": contentType,
	})
}
