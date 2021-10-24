package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Supervisor struct {
	m int32
	counter int32
	done chan struct{}
	err chan struct{}
}

func NewSupervisor(n, m int) *Supervisor  {
	s := Supervisor{m: int32(m)}

	s.done = make(chan struct{})
	s.err = make(chan struct{}, n)

	go s.ErrorsHandler()

	return &s
}

func (s *Supervisor) Close() {
	close(s.err)
}

func (s *Supervisor) Generator(tasks []Task) <-chan Task{
	ch := make(chan Task)

	go func() {
		defer close(ch)

		for _, task := range tasks {
			select {
			case ch <- task:
			case <-s.done:
				return
			}
		}
	}()

	return ch
}

func (s *Supervisor) Worker(tasks <-chan Task) {
	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}

			if err := task(); err != nil {
				s.err <- struct{}{}
			}
		case <-s.done:
			return
		}
	}
}

func (s *Supervisor) ErrorsHandler() {
	for  {
		select {
		case <-s.err:
			atomic.AddInt32(&s.counter, 1)

			if s.IsErrorsLimitExceeded() {
				close(s.done)
				return
			}
		case <-s.done:
			return
		}
	}
}

func (s *Supervisor) IsErrorsLimitExceeded() bool {
	return atomic.LoadInt32(&s.counter) >= s.m
}


// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	s := NewSupervisor(n, m)
	defer s.Close()

	generator := s.Generator(tasks)

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			s.Worker(generator)
		}()
	}

	wg.Wait()

	if s.IsErrorsLimitExceeded() {
		return ErrErrorsLimitExceeded;
	}

	return nil
}
