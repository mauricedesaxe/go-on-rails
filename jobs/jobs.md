# How Jobs Work Architecturally and Adding a New Job

## Architectural Overview

The job system in the web application is designed to handle tasks that can be executed asynchronously, either on a regular schedule or triggered by specific events. This system is built around a few key components:

1. **Job Interface**: All jobs implement the [Job](/jobs/runner.go) interface, which requires an [Execute()](/jobs/hello_job.go) method. This ensures consistency and allows for polymorphism when handling different types of jobs.

    
```go
type Job interface {
	Execute() error
}
```


2. **Scheduler**: The scheduler uses the `cron` package to run jobs at specific times. It's configured to execute tasks like daily jobs, and it can be extended to support various schedules.

    
```go
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
```


3. **Queue**: For irregular jobs that need to run based on events rather than time, a queue system is implemented. Jobs are added to the queue and executed asynchronously.

    
```go
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
```


4. **Async Execution**: Jobs can be run asynchronously, with support for retries and backoff in case of failures. This is crucial for ensuring that temporary issues don't prevent jobs from completing.

    
```go

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
```


## Adding a New Job

To add a new job to the application, follow these steps:

1. **Define the Job**: Create a new Go file for your job in the `jobs` package. Implement the `Job` interface by defining the `Execute()` method. This method should contain the logic you want to run.

    Example:
    ```go
    package jobs

    import "log"

    type MyNewJob struct {
        // Add any fields your job needs here
    }

    func (job MyNewJob) Execute() error {
        // Your job's logic here
        log.Println("Executing MyNewJob")
        return nil
    }
    ```

2. **Schedule the Job**: If your job needs to run on a schedule, add it to the `StartSchedule` function in `runner.go`. Use the `cron` syntax to specify when the job should run.

    Example:
    ```go
    _, err := c.AddFunc("@hourly", func() {
        myJob := MyNewJob{
            // Initialize your job's fields here
        }
        AsyncRun(myJob)
    })
    ```

3. **Add to Queue**: If your job should run based on events, use the `AddToQueue` function to add it to the queue whenever the triggering event occurs.

    Example:
    ```go
    myJob := MyNewJob{
        // Initialize your job's fields here
    }
    AddToQueue(myJob)
    ```

4. **Test**: Ensure your job executes as expected by writing tests that trigger its execution, either through the scheduler or the queue.

By following these steps, you can extend the application's functionality with background jobs that run on schedules or in response to events.