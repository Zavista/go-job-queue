package jobs

import "sync"

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
		wp.wg.Add(1)

		go func() { // each worker will read from job queue and process a job
			defer wp.wg.Done()

			for job := range wp.manager.jobQueue { // infinite loop waiting for job, worker will keep pulling/executing jobs until queue is closed
				wp.processJob(job)
			}
		}()
	}
}

package jobs

import "sync"

type WorkerPool struct {
	numWorkers int
	manager    *Manager
	wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int, manager *Manager) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		manager:    manager,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)

		go func() {
			defer wp.wg.Done()

			for job := range wp.manager.jobQueue {
				wp.processJob(job)
			}
		}()
	}
}

func (wp *WorkerPool) processJob(job *Job) {
	wp.manager.mu.Lock()
	job.Status = Running
	job.Attempts++
	wp.manager.mu.Unlock()

	result, err := job.Payload.Process()

	wp.manager.mu.Lock()
	defer wp.manager.mu.Unlock()

	if err != nil {
		job.Status = Failed
		job.Err = err
		return
	}

	job.Status = Success
	job.Result = result
}