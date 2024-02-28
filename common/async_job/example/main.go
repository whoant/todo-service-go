package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"todo-service/common/async_job"
)

func main() {
	job1 := async_job.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 2)
		fmt.Println("I am job 1")

		return nil
	}, async_job.WithName("Job 1"))

	job2 := async_job.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 3)
		fmt.Println("I am job 2")

		return nil
	}, async_job.WithName("Job 2"))

	if err := async_job.NewGroup(false, job1, job2).
		Run(context.Background()); err != nil {
		log.Println(err)
	}
}
