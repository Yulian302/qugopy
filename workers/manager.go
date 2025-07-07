package workers

import (
	"fmt"
	"sync"
)

type Worker interface {
	Start() error
	Stop() error
	HealthCheck() error
	ID() string
}

type WorkerManager struct {
	workers []Worker
	mu      sync.Mutex
}

func NewWorkerManager() *WorkerManager {
	return &WorkerManager{}
}

func (wm *WorkerManager) AddWorker(w Worker) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.workers = append(wm.workers, w)
}

func (wm *WorkerManager) StartAll() error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	for _, w := range wm.workers {
		if err := w.Start(); err != nil {
			return fmt.Errorf("failed to start worker %s: %w", w.ID(), err)
		}
	}
	return nil
}

func (wm *WorkerManager) StopAll() error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	var errs []error
	for _, w := range wm.workers {
		if err := w.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("worker %s: %w", w.ID(), err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors stopping workers: %v", errs)
	}
	return nil
}

func (wm *WorkerManager) HealthCheck() map[string]error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	status := make(map[string]error)
	for _, w := range wm.workers {
		status[w.ID()] = w.HealthCheck()
	}
	return status
}
