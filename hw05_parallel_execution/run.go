package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// ErrorCounter каунтер ошибок.
type ErrorCounter struct {
	mutex       sync.Mutex
	errorsCount int
}

// Add добавляем ошибку.
func (t *ErrorCounter) Add() int {
	defer t.mutex.Unlock()
	t.mutex.Lock()
	t.errorsCount++
	return t.errorsCount
}

// GetCount получаем текущее кол-во ошибок.
func (t *ErrorCounter) GetCount() int {
	defer t.mutex.Unlock()
	t.mutex.Lock()
	return t.errorsCount
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errorCounter := new(ErrorCounter)
	wg := sync.WaitGroup{}
	tasksChan := make(chan Task)

	// если задач меньше, чем n, меняем n
	if len(tasks) < n {
		n = len(tasks)
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		// запускаем воркеры
		go taskExecutor(tasksChan, &wg, errorCounter)
	}

	for _, task := range tasks {
		// проверяем кол-во ошибок
		if errorCounter.GetCount() >= m {
			break
		}
		// кладем задачу в канал
		tasksChan <- task
	}
	// закрываем канал
	close(tasksChan)
	// дожидаемся завершения работы воркеров
	wg.Wait()

	if errorCounter.GetCount() >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func taskExecutor(tasksChain <-chan Task, wg *sync.WaitGroup, errorCounter *ErrorCounter) {
	defer wg.Done()

	// выбираем задачу из канала, запускаем
	for task := range tasksChain {
		if err := task(); err != nil {
			// задача завершилась с ошибкой, фиксируем
			errorCounter.Add()
		}
	}
}
