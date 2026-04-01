package jobs

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/zavista/go-job-queue/internal/processors"
)

type Manager struct {
	mu       sync.Mutex   // allow only one goroutine to read/write to the shared manager state
	jobs     map[int]*Job // map a job's id to a Job
	nextID   int          // incremental counter to get a unique ID for next job (just increase by 1 each time)
	jobQueue chan *Job    // communication channel between manager/workers. manager adds job to channel, workers (goroutines) recieve from channel
	logger   *log.Logger
}

// NewManager initializes a Manager with valid values as nil map/channel causes panic when written to
func NewManager(buffer int) *Manager {
	return &Manager{
		jobs:     make(map[int]*Job),
		nextID:   1,
		jobQueue: make(chan *Job, buffer),
		logger:   log.New(os.Stdout, "[manager]", log.LstdFlags),
	}
}

// AddJob recieves a Processor, creates a Job for it, and adds it to the Job map and queue
func (m *Manager) AddJob(p processors.Processor) *Job {
	m.mu.Lock()

	job := &Job{
		ID:      m.nextID,
		Name:    fmt.Sprintf("%s-job-%d", p.Type(), m.nextID),
		Status:  Pending,
		Payload: p,
	}

	m.jobs[job.ID] = job
	m.nextID++

	m.mu.Unlock() // Must release mutex before sending to channel as that can block and does not require mutex

	m.logger.Printf("created job %d (%s), queueing now", job.ID, job.Name)

	m.jobQueue <- job

	m.logger.Printf("queued job %d (%s)", job.ID, job.Name)

	return job
}

func (m *Manager) GetJob(id int) (*Job, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, ok := m.jobs[id]
	return job, ok
}

func (m *Manager) ListJobs() []*Job {
	m.mu.Lock()
	defer m.mu.Unlock()

	jobList := make([]*Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		jobList = append(jobList, job)
	}

	return jobList
}

func (m *Manager) Logger() *log.Logger {
	return m.logger
}
