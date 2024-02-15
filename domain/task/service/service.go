package service

import (
	"context"
	"sync/atomic"
	"task/domain/task/model"
	"task/domain/task/storage"
	"task/internal/logging"
	"time"
)

var _ Service = &service{}

type Service interface {
	GetOrCreate(ctx context.Context, taskId int) model.Task
	GetAll(ctx context.Context) model.Tasks
	CountTasksInProgress(ctx context.Context) int32
	Shutdown(ctx context.Context) string
}

type service struct {
	logger  *logging.Logger
	storage storage.Storage

	progressTaskCounter atomic.Int32
}

func NewService(logger *logging.Logger, storage storage.Storage) Service {
	return &service{
		logger:  logger,
		storage: storage,
	}
}

func (s *service) StartTask(ctx context.Context, task model.Task) {
	s.progressTaskCounter.Add(1)
	s.logger.Info("start task work", "id", task.Id, "count", task.Count)
	time.Sleep(time.Millisecond * 100 * time.Duration(task.Count))

	s.progressTaskCounter.Add(-1)
	s.logger.Info("end task work", "id", task.Id, "count", task.Count)
}

func (s *service) GetOrCreate(ctx context.Context, taskId int) model.Task {
	task := s.storage.GetOrCreate(ctx, taskId)

	// task in progress
	s.StartTask(ctx, task)

	return task
}

func (s *service) GetAll(ctx context.Context) model.Tasks {
	return s.storage.GetAll(ctx)
}

func (s *service) CountTasksInProgress(ctx context.Context) int32 {
	return s.progressTaskCounter.Load()
}

func (s *service) Shutdown(ctx context.Context) string {
	tasks := s.GetAll(ctx)
	return tasks.String()
}
