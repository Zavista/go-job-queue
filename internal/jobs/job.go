package jobs

import (
	"log"
	"sync"

	"github.com/zavista/go-job-queue/internal/processors"
)

type Job struct {
	mu       sync.Mutex
	ID       int
	Name     string
	Status   JobStatus
	Payload  processors.Processor
	Result   string
	Err      error
	Attempts int
}

type JobSnapshot struct {
	ID       int
	Name     string
	Status   JobStatus
	Result   string
	Err      error
	Attempts int
}

func (j *Job) MarkRunning(logger *log.Logger) {
	j.mu.Lock()
	j.Status = Running
	j.Attempts++
	attempts := j.Attempts
	j.mu.Unlock()

	if logger != nil {
		logger.Printf("job %d (%s): status=%s attempts=%d", j.ID, j.Name, Running, attempts)
	}
}

func (j *Job) MarkFailed(err error, logger *log.Logger) {
	j.mu.Lock()
	j.Status = Failed
	j.Err = err
	j.mu.Unlock()

	if logger != nil {
		logger.Printf("job %d (%s): status=%s error=%v", j.ID, j.Name, Failed, err)
	}
}

func (j *Job) MarkSuccess(result string, logger *log.Logger) {
	j.mu.Lock()
	j.Status = Success
	j.Result = result
	j.mu.Unlock()

	if logger != nil {
		logger.Printf("job %d (%s): status=%s result=%q", j.ID, j.Name, Success, result)
	}
}

func (j *Job) Snapshot() JobSnapshot {
	j.mu.Lock()
	defer j.mu.Unlock()

	return JobSnapshot{
		ID:       j.ID,
		Name:     j.Name,
		Status:   j.Status,
		Result:   j.Result,
		Err:      j.Err,
		Attempts: j.Attempts,
	}
}
