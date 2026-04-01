package jobs

import (
	"fmt"
	"sync"

	"github.com/zavista/go-job-queue/internal/processors"
)

type Manager struct {
	mu       sync.Mutex   // allow only one goroutine to read/write to the shared manager state
	jobs     map[int]*Job // map a job's id to a Job
	nextID   int          // incremental counter to get a unique ID for next job (just increase by 1 each time)
	jobQueue chan *Job    // communication channel between manager/workers. manager adds job to channel, workers (goroutines) recieve from channel
}

// NewManager initializes a Manager with valid values as nil map/channel causes panic when written to
func NewManager(buffer int) *Manager {
	return &Manager{
		jobs:     make(map[int]*Job),
		nextID:   1,
		jobQueue: make(chan *Job, buffer),
	}
}

// AddJob recieves a Processor, creates a Job for it, and adds it to the Job map and queue
func (m *Manager) AddJob(p processors.Processor) *Job {
	m.mu.Lock()

	job := &Job{
		ID:      m.nextID,
		Name:    fmt.Sprintf("%v: job%v", p.Type(), m.nextID),
		Status:  Pending,
		Payload: p,
	}

	m.jobs[job.ID] = job
	m.nextID++

	m.mu.Unlock() // Must release mutex before sending to channel as that can block and does not require mutex

	m.jobQueue <- job
	return job
}
