package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/half0wl/railway-temporal/example/cron"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	TEMPORAL_HOST := os.Getenv("TEMPORAL_HOST")
	if TEMPORAL_HOST == "" {
		panic("TEMPORAL_HOST is not set")
	}

	fmt.Println("Temporal host: ", TEMPORAL_HOST)

	c, err := client.Dial(client.Options{
		HostPort: TEMPORAL_HOST,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	cmd := os.Args[1]
	if cmd == "worker" {
		fmt.Println("Starting worker")
		w := worker.New(c, "cron", worker.Options{})

		w.RegisterWorkflow(cron.SampleCronWorkflow)
		w.RegisterActivity(cron.DoSomething)

		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalln("Unable to start worker", err)
		}
	} else if cmd == "cron" {
		fmt.Println("Starting cron")
		workflowID := "cron_" + uuid.New().String()
		workflowOptions := client.StartWorkflowOptions{
			ID:           workflowID,
			TaskQueue:    "cron",
			CronSchedule: "* * * * *",
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, cron.SampleCronWorkflow)
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
	} else {
		panic("Invalid command")
	}

}
