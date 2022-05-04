package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

// CronJob Represent a mission
type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time //expr.Next(time.Now())
}

func main() {
	//A schedule coroutine that schedules jobs if they expire

	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression
		now           time.Time
		scheduleTable map[string]*CronJob // key:The name of job;
	)
	now = time.Now()
	scheduleTable = make(map[string]*CronJob)
	//Define 2 CronJob
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}

	//Jobs was registered in the scheduleTable
	scheduleTable["job1"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}

	//Jobs was registered in the scheduleTable
	scheduleTable["job2"] = cronJob

	//Start a schedule coroutine

	go func() {
		var (
			jobName string
			job     *CronJob
			now     time.Time
		)
		// Periodically check the task scheduling table
		for {
			now = time.Now()
			for jobName, job = range scheduleTable {
				//Determine expiration date
				if job.nextTime.Before(now) || job.nextTime.Equal(now) {
					//Start a coroutine
					go func(jobName string) {
						fmt.Println("execute:" + jobName)
					}(jobName)
				}
				//Calculating the next scheduling time
				job.nextTime = job.expr.Next(now)
				fmt.Println(jobName, "the next execution time is:", job.nextTime, "| now time is:", now)
			}

			select {
			case <-time.NewTimer(1 * time.Second).C: //Readable after 100ms, return

			}

			//time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(100 * time.Second)
}
