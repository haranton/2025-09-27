package worker

import "webfilehosting/internal/service"

type Dispatcher struct {
	taskService     *service.TaskService
	downloadService *service.DownloadService
	workerPool      *Pool
	stopChan        chan struct{}
}

func (d *Dispatcher) Start()
func (d *Dispatcher) Stop()
func (d *Dispatcher) recoverInterruptedTasks() // Восстановление при запуске!
