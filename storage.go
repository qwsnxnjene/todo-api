package main

import "sync"

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type Storage struct {
	mu     sync.RWMutex
	tasks  []Task
	nextID int
}

func NewStorage() *Storage {
	return &Storage{
		tasks:  make([]Task, 0),
		nextID: 1,
	}
}

func (s *Storage) Add(text string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := Task{ID: s.nextID, Text: text, Done: false}
	s.tasks = append(s.tasks, task)
	s.nextID++
	return task
}

func (s *Storage) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// возвращаем копию, чтобы не было гонок при изменении
	copyTasks := make([]Task, len(s.tasks))
	copy(copyTasks, s.tasks)
	return copyTasks
}
