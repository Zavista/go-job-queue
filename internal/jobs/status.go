package jobs

type JobStatus string

// Package-level, named constants so we can get ENUM-like behaviour
const (
	Pending JobStatus = "pending"
	Running JobStatus = "running"
	Success JobStatus = "success"
	Failed  JobStatus = "failed"
)
