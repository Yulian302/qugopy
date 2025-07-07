package workers

import (
	"context"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/Yulian302/qugopy/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkerLifecycle(t *testing.T) {
	t.Run("SuccessfulWorkerCreation", func(t *testing.T) {
		wm := NewWorkerManager()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		config := PythonWorkerConfig{
			EnvPath:      path.Join(config.ProjectRootPath, "processing", "venv", "bin"),
			FilePath:     path.Join(config.ProjectRootPath, "processing", "worker.py"),
			Mode:         "test",
			IsProduction: false,
		}

		// test worker addition
		worker := NewPythonWorker(ctx, uuid.New().String(), config)
		wm.AddWorker(worker)
		assert.Equal(t, 1, len(wm.workers), "Worker should be added to manager")

		// test worker startup
		err := wm.StartAll()
		require.NoError(t, err, "Workers should start successfully")
		assert.NoError(t, worker.HealthCheck(), "Worker should report running state")

		// test individual worker shutdown
		err = worker.Stop()
		require.NoError(t, err, "Worker should stop cleanly")
		assert.NoError(t, worker.HealthCheck(), "Worker should report stopped state")
	})

	t.Run("ConcurrentWorkerManagement", func(t *testing.T) {
		wm := NewWorkerManager()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var wg sync.WaitGroup
		workerCount := 5

		// concurrently add workers
		for i := 0; i < workerCount; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				config := PythonWorkerConfig{
					EnvPath:      path.Join(config.ProjectRootPath, "processing", "venv", "bin"),
					FilePath:     path.Join(config.ProjectRootPath, "processing", "worker.py"),
					Mode:         "test",
					IsProduction: false,
				}
				wm.AddWorker(NewPythonWorker(ctx, uuid.New().String(), config))
			}(i)
		}
		wg.Wait()

		assert.Equal(t, workerCount, len(wm.workers), "All workers should be added without race conditions")

		// test bulk startup
		err := wm.StartAll()
		require.NoError(t, err, "All workers should start concurrently")

		// verify all workers are running
		runningCount := 0
		for _, w := range wm.workers {
			if err := w.HealthCheck(); err == nil {
				runningCount++
			}
		}
		assert.Equal(t, workerCount, runningCount, "All workers should be running")
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		// test invalid worker config
		t.Run("InvalidPythonPath", func(t *testing.T) {
			wm := NewWorkerManager()
			ctx := context.Background()

			invalidConfig := PythonWorkerConfig{
				EnvPath:      "/invalid/path",
				FilePath:     "/invalid/script.py",
				Mode:         "test",
				IsProduction: false,
			}

			worker := NewPythonWorker(ctx, uuid.New().String(), invalidConfig)
			wm.AddWorker(worker)

			err := wm.StartAll()
			require.Error(t, err, "Should fail with invalid Python path")
			assert.Error(t, worker.HealthCheck(), "Worker should not be running after failed start")
		})
	})

	t.Run("GracefulShutdown", func(t *testing.T) {
		wm := NewWorkerManager()
		ctx, cancel := context.WithCancel(context.Background())

		config := PythonWorkerConfig{
			EnvPath:      path.Join(config.ProjectRootPath, "processing", "venv", "bin"),
			FilePath:     path.Join(config.ProjectRootPath, "processing", "worker.py"),
			Mode:         "test",
			IsProduction: false,
		}

		// add and start workers
		for i := 0; i < 3; i++ {
			wm.AddWorker(NewPythonWorker(ctx, uuid.New().String(), config))
		}
		require.NoError(t, wm.StartAll())

		// test context cancellation
		cancel()
		time.Sleep(100 * time.Millisecond) // allow for graceful shutdown

		// verify all workers stopped
		for _, w := range wm.workers {
			assert.NoError(t, w.HealthCheck(), "All workers should stop after context cancellation")
		}
	})
}
