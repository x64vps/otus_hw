package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Supervisor struct {
	m      int32
	errCnt int32
}

func NewSupervisor(m int) *Supervisor {
	return &Supervisor{m: int32(m)}
}

func (s *Supervisor) Generator(tasks []Task) <-chan Task {
	ch := make(chan Task)

	go func() {
		defer close(ch)

		for _, task := range tasks {
			if s.IsErrorsLimitExceeded() {
				return
			}

			ch <- task
		}
	}()

	return ch
}

func (s *Supervisor) IsErrorsLimitExceeded() bool {
	return atomic.LoadInt32(&s.errCnt) >= s.m
}

func (s *Supervisor) IncrementErrorsCounter()  {
	atomic.AddInt32(&s.errCnt, 1)
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n < 1 || m < 1 || len(tasks) == 0{
		return nil
	}

	s := NewSupervisor(m)

	generator := s.Generator(tasks)

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for task := range generator {
				if err := task(); err != nil {
					s.IncrementErrorsCounter()
				}
			}
		}()
	}

	wg.Wait()

	if s.IsErrorsLimitExceeded() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
