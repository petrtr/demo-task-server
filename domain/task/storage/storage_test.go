package storage

import (
	"context"
	"strings"
	"sync"
	"task/domain/task/model"
	"task/internal/logging"
	"testing"
)

func TestNewStorage(t *testing.T) {

	logger := logging.GetLogger("test", "info")
	storage := NewStorage(logger)

	if storage == nil {
		t.Errorf("storage should not be nil")
	}
}

func TestStorage_GetOrCreate(t *testing.T) {

	logger := logging.GetLogger("test", "info")
	storage := NewStorage(logger)

	task := storage.GetOrCreate(context.Background(), 1)

	if task.Id != 1 {
		t.Errorf("task id should be 1")
	}

	if task.Count != 1 {
		t.Errorf("task count should be 1")
	}

	countTestTask := 1000
	countRequest := 1000

	wg := &sync.WaitGroup{}

	for j := 2; j < countTestTask+2; j++ {
		j := j
		for i := 0; i < countRequest; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = storage.GetOrCreate(context.Background(), j)
			}()

		}
	}

	wg.Wait()

	for j := 2; j < countTestTask+2; j++ {
		task = storage.GetOrCreate(context.Background(), j)
		if task.Id != j {
			t.Errorf("task id should be j=%d", j)
		}
		if task.Count != uint(countRequest+1) {
			t.Errorf("task count for task id=%d should be %d", j, countRequest+1)
		}
	}

}

func TestStorage_GetAll(t *testing.T) {

	logger := logging.GetLogger("test", "info")
	storage := NewStorage(logger)

	tasks := storage.GetAll(context.Background())
	if len(tasks) != 0 {
		t.Errorf("tasks count should be %d", 0)
	}

	if tasks.String() != "" {
		t.Errorf("tasks string should be empty")
	}

	countTestTask := 1000
	countRequest := 1000

	wg := &sync.WaitGroup{}

	for j := 1; j < countTestTask+1; j++ {
		j := j
		for i := 0; i < countRequest; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = storage.GetOrCreate(context.Background(), j)
			}()
		}
	}

	wg.Wait()

	tasks = storage.GetAll(context.Background())

	if len(tasks) != countTestTask {
		t.Errorf("tasks count should be %d", countTestTask)
	}

	for _, task := range tasks {
		if task.Id < 1 || task.Id > countTestTask {
			t.Errorf("task id should be in range 1-%d", countTestTask)
		}
		if task.Count != uint(countRequest) {
			t.Errorf("task count for task id=%d should be %d", task.Id, countRequest)
		}
	}

	tasksFinalString := ""
	for i := 0; i < countTestTask; i++ {
		task := model.Task{
			Id:    i + 1,
			Count: uint(countRequest),
		}
		tasksFinalString += task.String() + "\n"
	}
	tasksFinalString = strings.Trim(tasksFinalString, "\n")

	if tasksFinalString == "" {
		t.Errorf("tasks string should not be empty")
	}

	if tasksFinalString != tasks.String() {
		t.Errorf("tasks string should be equal \n%s != \n%s", tasksFinalString, tasks.String())
	}

}
