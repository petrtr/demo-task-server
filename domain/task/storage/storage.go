package storage

import (
	"context"
	"sync"
	"task/domain/task/model"
	"task/internal/logging"
)

var _ Storage = &storage{}

type Storage interface {
	GetOrCreate(ctx context.Context, taskId int) model.Task
	GetAll(ctx context.Context) model.Tasks
}

type storage struct {
	logger *logging.Logger
	db     map[int]uint
	mx     *sync.Mutex
}

func NewStorage(logger *logging.Logger) Storage {

	return &storage{
		logger: logger,
		db:     make(map[int]uint),
		mx:     &sync.Mutex{},
	}
}

func (s *storage) GetOrCreate(ctx context.Context, taskId int) model.Task {

	s.mx.Lock()
	defer s.mx.Unlock()

	count, ok := s.db[taskId]
	if !ok {
		s.db[taskId] = 1
		return model.Task{
			Id:    taskId,
			Count: 1,
		}
	}

	// increment
	s.db[taskId] = count + 1

	return model.Task{
		Id:    taskId,
		Count: count + 1,
	}
}
func (s *storage) GetAll(ctx context.Context) model.Tasks {

	s.mx.Lock()
	defer s.mx.Unlock()

	tasks := make([]model.Task, 0, len(s.db))
	for id, count := range s.db {
		tasks = append(tasks, model.Task{
			Id:    id,
			Count: count,
		})
	}
	return tasks
}
