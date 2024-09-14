package main

import (
	"fmt"

	"github.com/labasubagia/task-tracker/cmd"
	"github.com/labasubagia/task-tracker/internal/task"
)

func main() {
	repo := task.NewRepo("task.json")
	svc := task.NewService(repo)
	err := cmd.Cmd(svc)
	if err != nil {
		fmt.Println("error", err)
	}
}
