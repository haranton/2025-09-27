package worker

import (
	"log"
	"sync"
	"time"
	"webfilehosting/internal/models"
	"webfilehosting/internal/service"
)

type Dispatcher struct {
	taskService     *service.TaskService
	downloadService *service.DownloadService
	workerPool      *Pool
	pollInterval    time.Duration
	stopChan        chan struct{}
	wg              sync.WaitGroup
}

func NewDispatcher(
	taskService *service.TaskService,
	downloadService *service.DownloadService,
	maxWorkers int,
) *Dispatcher {
	return &Dispatcher{
		taskService:     taskService,
		downloadService: downloadService,
		workerPool:      NewPool(maxWorkers),
		pollInterval:    5 * time.Second,
		stopChan:        make(chan struct{}),
	}
}

func (d *Dispatcher) Start() {
	log.Println("Starting dispatcher...")

	d.recoverInterruptedTasks()

	d.wg.Add(1)
	go d.run()

	log.Println("Dispatcher started")
}

func (d *Dispatcher) Stop() {
	log.Println("Stopping dispatcher...")

	close(d.stopChan)

	d.wg.Wait()

	d.workerPool.Stop()

	log.Println("Dispatcher stopped")
}

func (d *Dispatcher) run() {
	defer d.wg.Done()

	ticker := time.NewTicker(d.pollInterval)
	defer ticker.Stop()

	log.Println("Dispatcher main loop started")

	for {
		select {
		case <-d.stopChan:
			log.Println("Dispatcher main loop stopping")
			return

		case <-ticker.C:
			d.dispatchPendingTasks()
		}
	}
}

func (d *Dispatcher) recoverInterruptedTasks() {
	log.Println("Recovering interrupted tasks...")

	tasks, err := d.taskService.GetAllTasks()
	if err != nil {
		log.Printf("Error getting tasks for recovery: %v\n", err)
		return
	}

	recoveredCount := 0
	for _, task := range tasks {
		if task.Status == models.TaskStatusInProgress {
			log.Printf("Found interrupted task: %s\n", task.Id)

			// Сбрасываем статус задачи
			err := d.taskService.UpdateTaskStatus(task.Id.String(), models.TaskStatusPending)
			if err != nil {
				log.Printf("Error resetting task %s: %v\n", task.Id, err)
				continue
			}

			// Сбрасываем статусы файлов
			for i := range task.Urls {
				if task.Urls[i].Status == models.FileStatusDownloading {
					task.Urls[i].Status = models.FileStatusPending
					task.Urls[i].Error = ""
				}
			}

			// Сохраняем обновленную задачу
			err = d.taskService.SaveTask(task)
			if err != nil {
				log.Printf("Error saving recovered task %s: %v\n", task.Id, err)
				continue
			}

			recoveredCount++
			log.Printf("Task %s recovered and reset to pending\n", task.Id)
		}
	}

	log.Printf("Recovery completed: %d tasks recovered\n", recoveredCount)
}

func (d *Dispatcher) dispatchPendingTasks() {
	tasks, err := d.taskService.GetPendingTasks()
	if err != nil {
		log.Printf("Error getting pending tasks: %v\n", err)
		return
	}

	if len(tasks) > 0 {
		log.Printf("Found %d pending tasks\n", len(tasks))
	}

	for _, task := range tasks {
		select {
		case <-d.stopChan:
			return
		default:
			d.dispatchTask(task)
		}
	}
}

// dispatchTask отправляет одну задачу на выполнение
func (d *Dispatcher) dispatchTask(task *models.Task) {
	// Создаем процессор для задачи
	processor := NewTaskProcessor(d.taskService, d.downloadService, task)

	// Отправляем задачу в worker pool
	d.workerPool.Submit(processor.Process)

	log.Printf("Task %s dispatched to worker pool\n", task.Id)
}
