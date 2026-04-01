package jobs

import (
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
