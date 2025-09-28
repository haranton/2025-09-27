package worker

import (
	"fmt"
	"sync"
)

type Pool struct {
	jobs       chan func()
	wg         sync.WaitGroup
	stopChan   chan struct{}
	maxWorkers int
}

func NewPool(maxWorkers int) *Pool {
	if maxWorkers <= 0 {
		maxWorkers = 4
	}

	pool := &Pool{
		jobs:       make(chan func(), 100),
		stopChan:   make(chan struct{}),
		maxWorkers: maxWorkers,
	}

	pool.startWorkers()

	return pool

}

func (p *Pool) startWorkers() {
	for i := 0; i < p.maxWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

func (p *Pool) worker(workerID int) {
	defer p.wg.Done()
	for {
		select {
		case job, ok := <-p.jobs:
			if !ok {
				fmt.Printf("Worker %d stopping (channel closed)\n", workerID)
				return
			}
			fmt.Printf("Worker %d executing job\n", workerID)
			job()
			fmt.Printf("Worker %d job completed\n", workerID)
		case <-p.stopChan:
			fmt.Printf("Worker %d stopping (signal received)\n", workerID)
			return
		}
	}
}

func (p *Pool) Submit(job func()) {
	select {
	case p.jobs <- job:
		// Задача успешно добавлена в канал
		fmt.Println("Job submitted to worker pool")
	case <-p.stopChan:
		// Пул остановлен, задача не принимается
		fmt.Println("Worker pool is stopped, job rejected")
	default:
		// Канал заполнен, задача не может быть добавлена сразу
		fmt.Println("Worker pool is busy, job queued for retry")
	}
}

func (p *Pool) Stop() {

	fmt.Println("Stopping worker pool...")

	// Закрываем stopChan чтобы сигнализировать воркерам о остановке
	close(p.stopChan)

	// Ждем завершения всех воркеров
	p.wg.Wait()

	// Закрываем канал jobs
	close(p.jobs)

	fmt.Println("Worker pool stopped")
}

func (p *Pool) GetQueueSize() int {
	return len(p.jobs)
}
