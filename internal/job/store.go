package job

import (
	"fmt"
	"go-video-converter/internal/transcoder"
	"sync"
)

type JobStatus string

const (
	JobQueued    JobStatus = "queued"
	JobRunning   JobStatus = "running"
	JobCompleted JobStatus = "completed"
	JobFailed    JobStatus = "failed"
)

type Job struct {
	ID       string              `json:"id"`
	Status   JobStatus           `json:"status"`
	Input    string              `json:"input"`
	Output   string              `json:"output"`
	Progress transcoder.Progress `json:"progress"`
}

type Store struct {
	mu   sync.RWMutex
	jobs map[string]Job
}

func NewStore() *Store {
	store := Store{
		jobs: make(map[string]Job),
	}

	return &store
}

func (s *Store) Create(job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.jobs[job.ID]

	if exists {
		return fmt.Errorf("job with id already exists")
	}

	s.jobs[job.ID] = job
	return nil

}


func (s *Store) Get(id string) (Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[id]

	return job, exists
}


func (s *Store) Update(id string, updateFn func(*Job)) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.jobs[id]

	if !exists {
		return fmt.Errorf("job with id doesn't exist")
	}

	updateFn(&job)

	s.jobs[id] = job

	return nil
}
