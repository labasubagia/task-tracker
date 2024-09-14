package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/labasubagia/task-tracker/internal/task"
)

func Cmd(svc *task.Service) error {
	args := make([]string, 3)
	copy(args, os.Args[1:])

	switch args[0] {

	case "add":
		return svc.Add(task.NewTask(0, args[1]))

	case "update":
		ID, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("ID is not number")
		}
		return svc.UpdateDesc(int64(ID), args[2])

	case "mark-in-progress":
		ID, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("ID is not number")
		}
		return svc.MarkStatus(int64(ID), task.StatusInProgress)

	case "mark-done":
		ID, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("ID is not number")
		}
		return svc.MarkStatus(int64(ID), task.StatusDone)

	case "delete":
		ID, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("ID is not number")
		}
		return svc.Delete(int64(ID))

	case "list":
		tasks, err := svc.List(task.Status(args[1]))
		if err != nil {
			return err
		}
		for _, task := range tasks {
			fmt.Printf("Task ID: %d, Status: %s, Desc: %s\n", task.ID, task.Status, task.Description)
		}
	default:
		return fmt.Errorf("command '%s' not found", args[0])
	}

	return nil
}
