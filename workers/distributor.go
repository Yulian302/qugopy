package workers

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/Yulian302/qugopy/config"
	"github.com/google/uuid"
)

type WorkerDistributor struct {
	pyManager *WorkerManager
	goManager *WorkerManager
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewWorkerDistributor() *WorkerDistributor {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerDistributor{
		pyManager: NewWorkerManager(),
		goManager: NewWorkerManager(),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (wd *WorkerDistributor) DistributeWorkers(totalWorkers int, mode string, isProduction bool) (context.CancelFunc, error) {
	if totalWorkers < 1 {
		return nil, fmt.Errorf("totalWorkers must be at least 1")
	}

	totalWorkers = min(totalWorkers, runtime.NumCPU())
	pyCount := totalWorkers / 2
	goCount := totalWorkers - pyCount

	config := PythonWorkerConfig{
		EnvPath:      path.Join(config.ProjectRootPath, "processing", "venv", "bin"),
		FilePath:     path.Join(config.ProjectRootPath, "processing", "worker.py"),
		Mode:         mode,
		IsProduction: isProduction,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < pyCount; i++ {
		wd.wg.Add(1)
		wd.pyManager.AddWorker(NewPythonWorker(ctx, uuid.New().String(), config))
	}

	for i := 0; i < goCount; i++ {
		wd.wg.Add(1)
		wd.goManager.AddWorker(NewGoWorker(
			uuid.New().String(),
			func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return nil
				default:
					// TODO
					return nil
				}
			},
		))
	}

	if err := wd.pyManager.StartAll(); err != nil {
		wd.cleanup()
		return nil, fmt.Errorf("python worker startup failed: %w", err)
	}
	if err := wd.goManager.StartAll(); err != nil {
		wd.cleanup()
		return nil, fmt.Errorf("go worker startup failed: %w", err)
	}

	return cancel, nil
}

func (wd *WorkerDistributor) Shutdown() error {
	wd.cancel()

	done := make(chan struct{})
	go func() {
		wd.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("shutdown timeout")
	}
}

func (wd *WorkerDistributor) cleanup() {
	wd.cancel()
	_ = wd.pyManager.StopAll()
	_ = wd.goManager.StopAll()
}
