package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	NotifyListener struct {
		MinReconnectInterval int `json:"minReconnectInterval"`
		MaxReconnectInterval int `json:"maxReconnectInterval"`
	}
	TaskRunner struct {
		MaxConcurrentTasks int
		TaskDeadline       int
		RetryAfter         int
	}
	Channels []string
}

func loadConfig(cfg *Config) {
	b, err := os.ReadFile("config.json")
	if err != nil {
		LStderr.Println(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(b, cfg)
	if err != nil {
		LStderr.Println(err.Error())
		os.Exit(1)
	}
}
