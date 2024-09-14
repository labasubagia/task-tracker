package task_test

import (
	"testing"
	"time"

	"github.com/labasubagia/task-tracker/internal/task"
	"github.com/stretchr/testify/require"
)

func TestRepoList(t *testing.T) {
	taskDone := []*task.Task{
		{ID: 1, Description: "buy an apple", Status: task.StatusDone},
	}

	taskInProgress := []*task.Task{
		{ID: 2, Description: "jogging 1 km", Status: task.StatusInProgress},
	}

	taskTodo := []*task.Task{
		{ID: 3, Description: "build CLI app", Status: task.StatusTodo},
		{ID: 4, Description: "take a nap", Status: task.StatusTodo},
	}

	data := make([]*task.Task, 0, len(taskDone)+len(taskInProgress)+len(taskTodo))
	data = append(data, taskDone...)
	data = append(data, taskInProgress...)
	data = append(data, taskTodo...)

	filename := task.RandomFileJSON()
	err := task.WriteToFile(filename, data)
	require.Nil(t, err)
	defer task.DeleteFile(filename)
	r := task.NewRepo(filename)

	got, err := r.List("")
	require.Nil(t, err)
	require.Equal(t, data, got)

	got, err = r.List(task.StatusDone)
	require.Nil(t, err)
	require.Equal(t, taskDone, got)

	got, err = r.List(task.StatusInProgress)
	require.Nil(t, err)
	require.Equal(t, taskInProgress, got)

	got, err = r.List(task.StatusTodo)
	require.Nil(t, err)
	require.Equal(t, taskTodo, got)
}

func TestRepoGet(t *testing.T) {

	data := []*task.Task{
		{ID: 1, Description: "buy an apple", Status: task.StatusDone},
		{ID: 2, Description: "jogging 1 km", Status: task.StatusInProgress},
		{ID: 3, Description: "build CLI app", Status: task.StatusTodo},
		{ID: 4, Description: "take a nap", Status: task.StatusTodo},
	}

	filename := task.RandomFileJSON()
	err := task.WriteToFile(filename, data)
	require.Nil(t, err)
	defer task.DeleteFile(filename)
	r := task.NewRepo(filename)

	got, err := r.Get(1)
	require.Nil(t, err)
	require.Equal(t, data[0], got)

	got, err = r.Get(4)
	require.Nil(t, err)
	require.Equal(t, data[3], got)

	_, err = r.Get(-12)
	require.NotNil(t, err)
}

func TestRepoAdd(t *testing.T) {
	filename := task.RandomFileJSON()
	r := task.NewRepo(filename)
	defer task.DeleteFile(filename)

	err := r.Add(&task.Task{Description: "buy a phone"})
	require.Nil(t, err)
	data, err := r.List("")
	require.Nil(t, err)
	require.Len(t, data, 1)
	require.Equal(t, int64(1), data[len(data)-1].ID)

	err = r.Add(&task.Task{Description: "buy a watch"})
	require.Nil(t, err)
	data, err = r.List("")
	require.Nil(t, err)
	require.Len(t, data, 2)
	require.Equal(t, int64(2), data[len(data)-1].ID)
}

func TestRepoUpdate(t *testing.T) {

	now := time.Now()

	data := []*task.Task{
		{ID: 1, Description: "buy an apple", Status: task.StatusDone, UpdatedAt: now},
		{ID: 2, Description: "jogging 1 km", Status: task.StatusInProgress, UpdatedAt: now},
	}
	filename := task.RandomFileJSON()
	err := task.WriteToFile(filename, data)
	require.Nil(t, err)
	defer task.DeleteFile(filename)

	r := task.NewRepo(filename)

	err = r.Update(-12, &task.Task{})
	require.NotNil(t, err)

	payload, _ := r.Get(1)
	before := *payload
	payload.Description = "buy a tomato"
	err = r.Update(1, payload)
	require.Nil(t, err)
	item, err := r.Get(1)
	require.Nil(t, err)
	require.Equal(t, "buy a tomato", item.Description)
	require.Greater(t, item.UpdatedAt, before.UpdatedAt)

	payload, _ = r.Get(2)
	before = *payload
	payload.Status = task.StatusDone
	err = r.Update(2, payload)
	require.Nil(t, err)
	item, err = r.Get(2)
	require.Nil(t, err)
	require.Equal(t, task.StatusDone, item.Status)
	require.Greater(t, item.UpdatedAt, before.UpdatedAt)
}

func TestRepoDelete(t *testing.T) {
	data := []*task.Task{
		{ID: 1, Description: "buy an apple", Status: task.StatusDone},
		{ID: 2, Description: "jogging 1 km", Status: task.StatusInProgress},
	}
	filename := task.RandomFileJSON()
	err := task.WriteToFile(filename, data)
	require.Nil(t, err)
	defer task.DeleteFile(filename)

	r := task.NewRepo(filename)

	err = r.Delete(-12)
	require.NotNil(t, err)

	err = r.Delete(1)
	require.Nil(t, err)

	err = r.Delete(2)
	require.Nil(t, err)

	data, err = r.List("")
	require.Nil(t, err)
	require.Len(t, data, 0)
}
