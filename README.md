# Go - Job Scheduler

A small concurrent job scheduler system written in Go.

## Project

The system accepts different kinds of jobs, wraps them in a common `Job` type, and sends them through a shared queue. A pool of workers consumes jobs from the queue concurrently and updates each job's status as it runs.

A job moves through a lifecycle like this:

```text
pending -> running -> success
pending -> running -> failed
```

## Architecture

### Processor

A Processor represents the actual work being done.

It is defined as an interface so that different job types can implement the same contract.

Example responsibilities:

- send an email
- generate a report
- simulate some task that may succeed or fail

### Job

A Job is the system's tracked record of work.

It stores:

- a unique ID
- a name
- a status
- the processor payload
- the result
- any error
- number of attempts

The Job owns its own mutable state and protects it with a mutex.

### Manager

The Manager is the central coordinator.

It is responsible for:

- assigning job IDs
- storing jobs in memory
- pushing jobs onto the queue
- exposing access to jobs

### WorkerPool

The WorkerPool starts a fixed number of worker goroutines.

Each worker:

- receives a job from the queue
- marks it as running
- executes the processor
- updates the final status

## Folder Structure

```
go-job-queue
├── go.mod
├── README.md
├── cmd/
│ └── app/
│ └── main.go
└── internal/
├── jobs/
│ ├── job.go
│ ├── manager.go
│ ├── status.go
│ └── worker.go
└── processors/
│ ├── processor.go
│ └── email.go
```
