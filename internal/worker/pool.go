package worker

import (
	"log"
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
				log.Printf("Worker %d stopping (channel closed)\n", workerID)
				return
			}
			log.Printf("Worker %d executing job\n", workerID)
			job()
			log.Printf("Worker %d job completed\n", workerID)
		case <-p.stopChan:
			log.Printf("Worker %d stopping (signal received)\n", workerID)
			return
		}
	}
}

func (p *Pool) Submit(job func()) {
	select {
	case p.jobs <- job:
		log.Println("Job submitted to worker pool")
	case <-p.stopChan:
		log.Println("Worker pool is stopped, job rejected")
	default:
		log.Println("Worker pool is busy, job queued for retry")
	}
}

func (p *Pool) Stop() {
	log.Println("Stopping worker pool...")

	close(p.stopChan)

	p.wg.Wait()

	close(p.jobs)

	log.Println("Worker pool stopped")
}

func (p *Pool) GetQueueSize() int {
	return len(p.jobs)
}
