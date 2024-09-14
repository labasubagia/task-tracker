package task

import (
	"errors"
	"time"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) Add(item *Task) error {
	now := time.Now()

	if item.Description == "" {
		return errors.New("please provide task description")
	}

	if item.Status == "" {
		item.Status = StatusTodo
	}

	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}

	if item.UpdatedAt.IsZero() {
		item.UpdatedAt = now
	}
	return s.repo.Add(item)
}

func (s *Service) List(status Status) ([]*Task, error) {
	isStatusValid := status == "" || s.isStatusValid(status)
	if !isStatusValid {
		return nil, errors.New("unsupported status")
	}

	return s.repo.List(status)
}

func (s *Service) Get(ID int64) (*Task, error) {
	return s.repo.Get(ID)
}

func (s *Service) UpdateDesc(ID int64, desc string) error {
	cur, err := s.repo.Get(ID)
	if err != nil {
		return err
	}
	cur.Description = desc
	return s.repo.Update(ID, cur)
}

func (s *Service) MarkStatus(ID int64, status Status) error {
	cur, err := s.repo.Get(ID)
	if err != nil {
		return err
	}
	if !s.isStatusValid(status) {
		return errors.New("unsupported status")
	}
	cur.Status = status
	return s.repo.Update(ID, cur)
}

func (s *Service) Delete(ID int64) error {
	return s.repo.Delete(ID)
}

func (s *Service) isStatusValid(status Status) bool {
	return status == StatusTodo ||
		status == StatusInProgress ||
		status == StatusDone
}
