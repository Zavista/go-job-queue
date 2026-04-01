package jobs

import "github.com/zavista/go-job-queue/internal/processors"

type Job struct {
	ID       int
	Name     string
	Status   JobStatus
	Payload  processors.Processor
	Result   string
	Err      error
	Attempts int
}
