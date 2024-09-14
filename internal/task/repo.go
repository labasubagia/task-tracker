package task

import (
	"fmt"
	"time"
)

type Repo struct {
	filename string
}

func NewRepo(filename string) *Repo {
	r := &Repo{
		filename: filename,
	}
	return r
}

func (r *Repo) List(status Status) ([]*Task, error) {
	data, err := ReadFromFile(r.filename)
	if err != nil {
		return nil, err
	}

	if status == "" {
		return data, nil
	}

	result := make([]*Task, 0)
	for _, item := range data {
		if item.Status == status {
			result = append(result, item)
		}
	}

	return result, nil
}

func (r *Repo) Get(ID int64) (*Task, error) {
	data, err := ReadFromFile(r.filename)
	if err != nil {
		return nil, err
	}

	index := -1
	for i, item := range data {
		if item.ID == ID {
			index = i
			break
		}
	}
	if index < 0 {
		return nil, fmt.Errorf("Task with ID: %d not found", ID)
	}

	return data[index], nil
}

func (r *Repo) Add(item *Task) error {
	data, err := ReadFromFile(r.filename)
	if err != nil {
		return err
	}

	if item.ID <= 0 {
		item.ID = 1
		if len(data) > 0 {
			item.ID = data[len(data)-1].ID + 1
		}
	}

	data = append(data, item)
	err = WriteToFile(r.filename, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(ID int64, item *Task) error {
	data, err := ReadFromFile(r.filename)
	if err != nil {
		return err
	}

	index := -1
	for i, item := range data {
		if item.ID == ID {
			index = i
			break
		}
	}
	if index < 0 {
		return fmt.Errorf("Task with ID: %d not found", ID)
	}

	item.UpdatedAt = time.Now()
	data[index] = item
	err = WriteToFile(r.filename, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Delete(ID int64) error {
	data, err := ReadFromFile(r.filename)
	if err != nil {
		return err
	}

	index := -1
	for i, item := range data {
		if item.ID == ID {
			index = i
			break
		}
	}
	if index < 0 {
		return fmt.Errorf("Task with ID: %d not found", ID)
	}

	data = append(data[:index], data[index+1:]...)
	err = WriteToFile(r.filename, data)
	if err != nil {
		return err
	}

	return nil
}
