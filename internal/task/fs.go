package task

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/google/uuid"
)

func getFile(filename string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filepath := path.Join(dir, filename)
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if info.Size() == 0 {
		_, err = file.WriteString("[]")
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

func DeleteFile(filename string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	filepath := path.Join(dir, filename)

	return os.Remove(filepath)
}

func ReadFromFile(filename string) ([]*Task, error) {
	file, err := getFile(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tasks []*Task
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, nil
	}

	err = json.Unmarshal(content, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func WriteToFile(filename string, tasks []*Task) error {
	file, err := getFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()

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
	return nil
}

func RandomFileJSON() string {
	return fmt.Sprintf("%s.json", uuid.NewString())
}
