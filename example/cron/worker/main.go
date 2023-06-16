package main

import (
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/cron"
)

func main() {
	TEMPORAL_HOST := os.Getenv("TEMPORAL_HOST")
	c, err := client.Dial(client.Options{
		HostPort: TEMPORAL_HOST,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "cron", worker.Options{})

	w.RegisterWorkflow(cron.SampleCronWorkflow)
	w.RegisterActivity(cron.DoSomething)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
