package main

import (
	"context"
	"crypto/rand"
	"sync"
	"time"
)

const (
	StatusEnded   = 0
	StatusOngoing = 1
)

type Scheduler struct {
	Ctx                context.Context
	MaxTasks           int
	TaskDeadline       time.Duration
	DefinedTasks       map[string]Task
	Ongoing            int
	Lock               sync.Mutex
	MaxConcurrentTasks int
	RetryAfter         time.Duration
}

func (sched *Scheduler) NewTask(payload, id string, task string) {
	if sched.Ongoing >= sched.MaxConcurrentTasks {
		go func() {
			time.Sleep(sched.RetryAfter)
			sched.NewTask(payload, id, task)
		}()
		LStderr.Printf("Max concurrent task count reached, waiting and retrying in %s ...", sched.RetryAfter.String())
		return
	}

	uuid := make([]byte, 15)
	_, err := rand.Read(uuid)
	if err != nil {
		LStderr.Println("error:", err)
	}
	status := make(chan int, 1)
	go sched.Track(uuid, &status)
	go sched.DefinedTasks[task](payload, id, &status)
}

func (sched *Scheduler) Track(uuid []byte, status *chan int) {
	for {
		select {
		case n := <-(*status):
			sched.Lock.Lock()
			if n == StatusOngoing {
				sched.Ongoing++
			} else if n == StatusEnded {
				sched.Ongoing--
				LStdout.Println("Received task finished - ending ongoing task track")
				return
			}
			sched.Lock.Unlock()
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}
