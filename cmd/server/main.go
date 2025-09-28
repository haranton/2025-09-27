package main

import (
	"fmt"
	"net/http"
	"os"
	"webfilehosting/internal/config"
	"webfilehosting/internal/service"
	"webfilehosting/internal/storage"
	"webfilehosting/internal/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализируем конфигурацию
	if err := config.InitConfig(); err != nil {
		fmt.Printf("Ошибка инициализации конфигурации: %v\n", err)
		os.Exit(1)
	}

	// Инициализация сервисов
	fileStorage := &storage.FileStorage{
		basePath:     config.AppConfig.StoragePath,
		tasksDir:     "tasks",
		downloadsDir: "downloads",
	}
	taskService := &service.TaskService{storage: fileStorage}
	downloadService := service.NewDownloadService()
	dispatcher := worker.NewDispatcher(taskService, downloadService, config.AppConfig.MaxWorkers)
	dispatcher.Start()
	defer dispatcher.Stop()

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

	r.POST("/tasks", func(ctx *gin.Context) {
		var body struct {
			Urls []string `json:"urls"`
		}
		if err := ctx.ShouldBindJSON(&body); err != nil || len(body.Urls) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON или пустой список URL"})
			return
		}
		task, err := taskService.CreateTask(body.Urls)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания задачи", "details": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"task_id": task.Id.String(), "status": task.Status})
	})

	r.GET("/tasks/:id", func(ctx *gin.Context) {
		taskID := ctx.Param("id")
		task, err := taskService.GetTask(taskID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
			return
		}
		ctx.JSON(http.StatusOK, task)
	})

	// Запускаем сервер с адресом из конфигурации
	serverAddr := config.AppConfig.GetServerAddress()
	fmt.Printf("Сервер запущен на %s\n", serverAddr)
	r.Run(serverAddr)
}
