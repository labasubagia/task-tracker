package task_test

import (
	"testing"

	"github.com/labasubagia/task-tracker/internal/task"
	"github.com/stretchr/testify/require"
)

func TestServiceAdd(t *testing.T) {
	filename := task.RandomFileJSON()
	defer task.DeleteFile(filename)
	repo := task.NewRepo(filename)
	svc := task.NewService(repo)

	item := task.NewTask(0, "buy new clothes")
	err := svc.Add(item)
	require.Nil(t, err)

	item = task.NewTask(0, "buy new phone")
	err = svc.Add(item)
	require.Nil(t, err)
}

func TestServiceList(t *testing.T) {
	filename := task.RandomFileJSON()
	defer task.DeleteFile(filename)
	repo := task.NewRepo(filename)
	svc := task.NewService(repo)

	payloads := []*task.Task{
		task.NewTask(1, "t1"),
		task.NewTask(2, "t2"),
		task.NewTask(3, "t3"),
		{ID: 4, Description: "t4", Status: task.StatusInProgress},
		{ID: 5, Description: "t5", Status: task.StatusDone},
	}
	for _, item := range payloads {
		err := svc.Add(item)
		require.Nil(t, err)
	}

	data, err := svc.List("")
	require.Nil(t, err)
	require.Len(t, data, len(payloads))

	data, err = svc.List(task.StatusTodo)
	require.Nil(t, err)
	require.Len(t, data, 3)

	data, err = svc.List(task.StatusInProgress)
	require.Nil(t, err)
	require.Len(t, data, 1)

	data, err = svc.List(task.StatusDone)
	require.Nil(t, err)
	require.Len(t, data, 1)
}

func TestServiceGet(t *testing.T) {
	filename := task.RandomFileJSON()
	defer task.DeleteFile(filename)
	repo := task.NewRepo(filename)
	svc := task.NewService(repo)

	data := []*task.Task{
		task.NewTask(1, "t1"),
		task.NewTask(2, "t2"),
	}
	for _, item := range data {
		err := svc.Add(item)
		require.Nil(t, err)
	}

	item, err := svc.Get(1)
	require.Nil(t, err)
	require.Equal(t, data[0].Description, item.Description)

	item, err = svc.Get(2)
	require.Nil(t, err)
	require.Equal(t, data[1].Description, item.Description)

	// fail
	item, err = svc.Get(-2)
	require.NotNil(t, err)
}

func TestServiceUpdateDesc(t *testing.T) {
	filename := task.RandomFileJSON()
	defer task.DeleteFile(filename)
	repo := task.NewRepo(filename)
	svc := task.NewService(repo)

	item := task.NewTask(1, "buy new clothes")
	err := svc.Add(item)
	require.Nil(t, err)

	err = svc.UpdateDesc(1, "buy new phone")
	require.Nil(t, err)

	result, err := svc.Get(1)
	require.Nil(t, err)
	require.Equal(t, item.ID, result.ID)
	require.Equal(t, "buy new phone", result.Description)
	require.Greater(t, result.UpdatedAt, item.UpdatedAt)

	err = svc.UpdateDesc(-100, "not found data")
	require.NotNil(t, err)
}

func TestServiceMarkStatus(t *testing.T) {
	filename := task.RandomFileJSON()
	defer task.DeleteFile(filename)
	repo := task.NewRepo(filename)
	svc := task.NewService(repo)

	item := task.NewTask(1, "buy new clothes")
	err := svc.Add(item)
	require.Nil(t, err)

	err = svc.MarkStatus(1, task.StatusInProgress)
	require.Nil(t, err)

	result, err := svc.Get(1)
	require.Nil(t, err)
	require.Equal(t, item.ID, result.ID)
	require.Equal(t, task.StatusInProgress, result.Status)
	require.Greater(t, result.UpdatedAt, item.UpdatedAt)

	// not found
	err = svc.MarkStatus(-200, task.StatusInProgress)
	require.NotNil(t, err)

	// invalid status
	err = svc.MarkStatus(1, task.Status("INVALID_STATUS"))
	require.NotNil(t, err)
}

func TestServiceDelete(t *testing.T) {
	filename := task.RandomFileJSON()
	defer task.DeleteFile(filename)
	repo := task.NewRepo(filename)
	svc := task.NewService(repo)

	item := task.NewTask(1, "buy new clothes")
	err := svc.Add(item)
	require.Nil(t, err)

	data, err := svc.List("")
	require.Nil(t, err)
	require.Len(t, data, 1)

	err = svc.Delete(item.ID)
	require.Nil(t, err)

	data, err = svc.List("")
	require.Nil(t, err)
	require.Len(t, data, 0)

	// fail
	err = svc.Delete(-21)
	require.NotNil(t, err)
}
