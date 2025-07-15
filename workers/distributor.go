package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/Yulian302/qugopy/config"
	"github.com/Yulian302/qugopy/internal/queue"
	"github.com/Yulian302/qugopy/internal/tasks"
	"github.com/Yulian302/qugopy/logging"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type WorkerDistributor struct {
	pyManager *WorkerManager
	goManager *WorkerManager
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	rdb       *redis.Client
}

func NewWorkerDistributor(rdb *redis.Client) *WorkerDistributor {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerDistributor{
		pyManager: NewWorkerManager(),
		goManager: NewWorkerManager(),
		ctx:       ctx,
		cancel:    cancel,
		rdb:       rdb,
	}
}

func (wd *WorkerDistributor) DistributeWorkers(totalWorkers int, mode string, isProduction bool, rdb *redis.Client) (context.CancelFunc, error) {
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

	ctx := wd.ctx

	for i := 0; i < pyCount; i++ {
		wd.wg.Add(1)
		wd.pyManager.AddWorker(NewPythonWorker(ctx, uuid.New().String(), config))
	}

	for i := 0; i < goCount; i++ {
		wd.wg.Add(1)
		wd.goManager.AddWorker(NewGoWorker(
			uuid.New().String(),
			func(ctx context.Context) error {
				for {
					select {
					case <-ctx.Done():
						return nil
					default:
						var task queue.IntTask
						var exists bool

						if mode == "redis" {
							res, err := rdb.ZPopMin("go_queue", 1).Result()
							if err != nil || len(res) == 0 {
								time.Sleep(100 * time.Millisecond)
								continue
							}
							if err := json.Unmarshal([]byte(res[0].Member.(string)), &task); err != nil {
								// skip task
								fmt.Printf("Failed to unmarshal task: %v. Raw: %s\n", err, res[0].Member)
								time.Sleep(100 * time.Millisecond)
								continue
							}

						} else {
							task, exists = queue.GoLocalQueue.PQ.Pop()
							// queue is empty
							if !exists {
								time.Sleep(100 * time.Millisecond)
								continue
							}
						}

						err := tasks.DispatchTask(ctx, task)
						if err != nil {
							logging.DebugLog(fmt.Sprintf("could not complete task (id=%s): %v", task.ID, err))
							continue
						}
					}
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

	return nil, nil
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
