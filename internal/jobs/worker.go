package jobs

import (
	"fmt"
	"sync"
)

type WorkerPool struct {
	numWorkers int
	manager    *Manager
	wg         sync.WaitGroup // tracks working go routines
}

func NewWorkerPool(numWorkers int, manager *Manager) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		manager:    manager,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		workerID := i + 1
		wp.wg.Add(1)

		go func(id int) { // each worker will be a goroutine
			defer wp.wg.Done()
			wp.manager.logger.Printf("worker %d started", id)

			for job := range wp.manager.jobQueue { // infinite loop waiting for job, worker will keep pulling/executing jobs until queue is closed
				wp.manager.logger.Printf("worker %d picked up job %d (%s)", id, job.ID, job.Name)
				wp.processJob(id, job)
			}

			wp.manager.logger.Printf("worker %d stopped", id)
		}(workerID)
	}
}

func (wp *WorkerPool) processJob(workerID int, job *Job) {
	job.MarkRunning(wp.manager.logger)

	result, err := job.Payload.Process()

	if err != nil {
		wrappedErr := fmt.Errorf("worker %d failed processing job %d: %w", workerID, job.ID, err)
		job.MarkFailed(wrappedErr, wp.manager.logger)
		return
	}

	job.MarkSuccess(result, wp.manager.logger)
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
