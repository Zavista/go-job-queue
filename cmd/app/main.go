package main

import (
	"fmt"

	"github.com/zavista/go-job-queue/internal/jobs"
	"github.com/zavista/go-job-queue/internal/processors"
)

func main() {
	manager := jobs.NewManager(10)
	workerPool := jobs.NewWorkerPool(2, manager)

	workerPool.Start()

	job1 := manager.AddJob(processors.EmailJob{
		To:      "a@test.com",
		Subject: "welcome",
	})

	job2 := manager.AddJob(processors.EmailJob{
		To:      "b@test.com",
		Subject: "hello",
	})

	job3 := manager.AddJob(processors.EmailJob{
		To:      "",
		Subject: "this will fail",
	})

	manager.CloseQueue()
	workerPool.Wait()

	fmt.Println()
	fmt.Println("final snapshots:")

	for _, job := range []*jobs.Job{job1, job2, job3} {
		snapshot := job.Snapshot()
		fmt.Printf(
			"job=%d name=%s status=%s attempts=%d result=%q err=%v\n",
			snapshot.ID,
			snapshot.Name,
			snapshot.Status,
			snapshot.Attempts,
			snapshot.Result,
			snapshot.Err,
		)
	}
}
