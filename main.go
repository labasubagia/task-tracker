package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func main() {
	err := RunTask()
	if err != nil {
		fmt.Println("error", err)
	}
}

type TaskStatus string

type Task struct {
	ID          int64      `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in-progress"
	StatusDone       TaskStatus = "done"
)

func RunTask() error {

	filename := "task.json"

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	var tasks []*Task

	if info.Size() > 0 {

		content, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		err = json.Unmarshal(content, &tasks)
		if err != nil {
			return err
		}
	}
	now := time.Now()

	args := make([]string, 3)
	copy(args, os.Args[1:])

	switch args[0] {

	case "add":
		desc := args[1]
		if desc == "" {
			return errors.New("please provide description")
		}
		var ID int64 = 1
		if len(tasks) > 0 {
			ID = tasks[len(tasks)-1].ID + 1
		}

		tasks = append(tasks, &Task{
			ID:          ID,
			Description: desc,
			Status:      StatusTodo,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		err = file.Truncate(0)
		if err != nil {
			return err
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(tasks)
		if err != nil {
			return err
		}
		_, err = file.Write(bytes)
		if err != nil {
			return err
		}

	case "update":
		ID, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("ID is not number")
		}
		desc := args[2]
		if desc == "" {
			return errors.New("please provide updated description")
		}
		if len(tasks) == 0 {
			return errors.New("there are no tasks yet")
		}
		found := false
		for _, task := range tasks {
			if task.ID == int64(ID) {
				found = true
				task.Description = desc
				task.UpdatedAt = now
				break
			}
		}
		if !found {
			return fmt.Errorf("task with id %d not found", ID)
		}

		err = file.Truncate(0)
		if err != nil {
			return err
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(tasks)
		if err != nil {
			return err
		}
		_, err = file.Write(bytes)
		if err != nil {
			return err
		}

	case "mark-in-progress":
		ID, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("ID is not number")
		}
		if len(tasks) == 0 {
			return errors.New("there are no tasks yet")
		}
		found := false
		for _, task := range tasks {
			if task.ID == int64(ID) {
				found = true
				task.Status = StatusInProgress
				task.UpdatedAt = now
				break
			}
		}
		if !found {
			return fmt.Errorf("task with id %d not found", ID)
		}

		err = file.Truncate(0)
		if err != nil {
			return err
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(tasks)
		if err != nil {
			return err
		}
		_, err = file.Write(bytes)
		if err != nil {
			return err
		}

	case "mark-done":
		ID, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("ID is not number")
		}
		if len(tasks) == 0 {
			return errors.New("there are no tasks yet")
		}
		found := false
		for _, task := range tasks {
			if task.ID == int64(ID) {
				found = true
				task.Status = StatusDone
				task.UpdatedAt = now
				break
			}
		}
		if !found {
			return fmt.Errorf("task with id %d not found", ID)
		}

		err = file.Truncate(0)
		if err != nil {
			return err
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(tasks)
		if err != nil {
			return err
		}
		_, err = file.Write(bytes)
		if err != nil {
			return err
		}

	case "delete":
		ID, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("ID is not number")
		}
		removeIndex := -1
		for i, task := range tasks {
			if task.ID == int64(ID) {
				removeIndex = i
			}
		}
		if removeIndex < 0 {
			return fmt.Errorf("task with id %d not found", ID)
		}

		tasks = append(tasks[:removeIndex], tasks[removeIndex+1:]...)

		err = file.Truncate(0)
		if err != nil {
			return err
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(tasks)
		if err != nil {
			return err
		}
		_, err = file.Write(bytes)
		if err != nil {
			return err
		}

	case "list":
		status := TaskStatus(args[1])
		isStatusValid := status == StatusTodo || status == StatusInProgress || status == StatusDone
		if status != "" && !isStatusValid {
			return errors.New("status unsupported")
		}

		for _, task := range tasks {
			isFiltered := status != "" && task.Status == status
			noFilter := status == ""
			if isFiltered || noFilter {
				fmt.Printf("Task ID: %d, Status: %s, Desc: %s\n", task.ID, task.Status, task.Description)
			}
		}

	default:
		fmt.Printf("command '%s' not found\n", args[0])
	}

	return nil
}
