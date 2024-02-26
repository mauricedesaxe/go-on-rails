// This is the core file of the jobs package as it defines the interfaces, global variables and helper functions
// that are used to run jobs both regularly (i.e. do X every 1 hour) and irregularly through a queue (i.e. do X
// when Y happens).

package jobs

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

// =================================================================================================
// Common job logic
// =================================================================================================

// Job is an interface that all jobs must implement
// It's meant to enforce a common method signature for all jobs and is usable for both regular and irregular jobs.
type Job interface {
	Execute() error
}

// AsyncRun attempts to execute a job asynchronously up to 3 times before giving up. It logs the time it
// took to execute the job and the number of attempts. It uses an exponential backoff strategy.
// This is a helper function that can be used to run any job asynchronously & is used to run both regular and irregular jobs.
func AsyncRun(job Job) {
	go func() {
		startTime := time.Now()
		var backoff = 1 * time.Second
		var attempts int
		for attempts = 0; attempts < 3; attempts++ {
			if err := job.Execute(); err != nil {
				log.Printf("Job %T failed to execute: %v", job, err)
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			break
		}
		log.Printf("Job %T took %v and %d attempts to execute", job, time.Since(startTime), attempts+1)
	}()
}

// sem is a semaphore to limit the number of concurrent jobs
var sem = make(chan struct{}, 10)

// AsyncRunLimited is a helper function that limits the number of concurrent jobs to the
// number of tokens in the semaphore. It's used to run irregular jobs.
// It's used to run irregular jobs.
func AsyncRunLimited(job Job) {
	sem <- struct{}{} // Acquire a token
	go func() {
		defer func() { <-sem }() // Release the token
		AsyncRun(job)
	}()
}

// =================================================================================================
// Irregularly scheduled jobs
// =================================================================================================

// queue is a channel for all irregularly scheduled jobs to be added to and processed by the runner.
// we chose to use a simple native go channel instead of other ways of managing the queue, 
// but you can also integrate with RabbitMQ/Kafka for this or even use a SQL or Redis db for this.
var queue chan Job = make(chan Job)

// shutdownChan is a channel used to signal the runner to stop processing irregular jobs
var shutdownChan = make(chan struct{})

// StartQueue is the central controller for all irregularly scheduled jobs
// This function is meant to be called once at the start of the application and it will run in the background.
// It listens for jobs to be added to the queue and runs them asynchronously when they are added.
// It also listens for a shutdown signal and stops processing jobs when it receives it.
func StartQueue() {
	for {
		select {
		case job := <-queue:
			// we use a limited async run to limit the number of concurrent jobs
			// this is done to reduce the load on the system and it's okay here
			// because irregular jobs are not time-sensitive and can be run at any time
			AsyncRunLimited(job)
		case <-shutdownChan:
			log.Println("Stopping irregular job processing...")
			return
		}
	}
}

// StopQueue signals the runner to stop processing irregular jobs
func StopQueue() {
	close(shutdownChan)
}

// AddToQueue adds a job to the Queue
func AddToQueue(job Job) {
	queue <- job
}

// =================================================================================================
// Regularly scheduled jobs
// =================================================================================================

// StartSchedule is the central controller for all regularly scheduled jobs in the application
// This one is manually configured to run jobs at specific times. It uses the cron package to do this.
// This function is meant to be called once at the start of the application and it will run in the background.
func StartSchedule() {
	go func() {
		c := cron.New()

		_, err := c.AddFunc("@daily", func() {
			log.Println("Executing daily jobs...")

			helloJob := HelloJob{
				Greeting: "Good morning",
				Name:     "John",
			}

			// we use an unlimited async run to run the job as soon as possible
			// this is because cron scheduled jobs are supposed to be run at specific times
			AsyncRun(helloJob)
		})
		if err != nil {
			log.Printf("Failed to schedule daily job: %v", err)
		}

		// Add more scheduled jobs here

		c.Start()

		// Graceful shutdown
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		c.Stop()
		log.Println("Job scheduler stopped")
	}()
}
