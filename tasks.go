package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type DemoTaskPayload struct {
	PrintThis string
}

const (
	StatusEnded   = 0
	StatusOngoing = 1
)

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
