/*
*	Add handlers for tasks in this file
*   - A handler should be of type Task
*	- Should properly flag when it started running and when it finishes via the constants `StatusOngoing` and `StatusEnded`                                     `
 */
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Task func(payload string, eventId string, status *chan int)

type DemoTaskPayload struct {
	PrintThis string
}

// MapTask add your specific task to a channel
func MapTask() {
	sched.DefinedTasks["Channel01"] = DemoTask // We map DemoTask defined below to pg_notify Channel01
}

/*
*	Declare new tasks of type Task below
 */

func DemoTask(payload string, eventId string, status *chan int) {
	(*status) <- StatusOngoing
	now := time.Now()
	j := DemoTaskPayload{}
	if err := json.Unmarshal([]byte(payload), &j); err != nil {
		LStderr.Println(err.Error())
	}
	time.Sleep(time.Second * 3) // processing delay
	LStdout.Println(j.PrintThis)
	t := time.Since(now).Seconds()
	if _, err := DB.Exec(fmt.Sprintf("select set_completed(%s,%f);", eventId, t)); err != nil {
		LStderr.Println(err.Error())
	}
	LStdout.Printf("DemoTask finished - eventId: %s - duration %fs\n", eventId, t)
	(*status) <- StatusEnded
}
