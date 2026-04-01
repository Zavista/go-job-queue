package main

import (
	"fmt"
	"time"

	"github.com/zavista/go-job-queue/internal/jobs"
	"github.com/zavista/go-job-queue/internal/processors"
)

func main() {
	manager := jobs.NewManager(10)
	workerPool := jobs.NewWorkerPool(2, manager)

	workerPool.Start()

	job := manager.AddJob(processors.EmailJob{
		To:      "test@example.com",
		Subject: "welcome",
	})

	fmt.Println("added job:", job.ID, job.Name, job.Status)

	time.Sleep(3 * time.Second)

	fmt.Println("final job state:")
	fmt.Println("id:", job.ID)
	fmt.Println("name:", job.Name)
	fmt.Println("status:", job.Status)
	fmt.Println("attempts:", job.Attempts)
	fmt.Println("result:", job.Result)

	if job.Err != nil {
		fmt.Println("error:", job.Err)
	}

}
