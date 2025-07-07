package workers

import (
	"errors"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockWorker implements Worker interface for testing
type MockWorker struct {
	id        string
	startErr  error
	stopErr   error
	healthErr error
}

func (m *MockWorker) Start() error { return m.startErr }
func (m *MockWorker) Stop() error  { return m.stopErr }
func (m *MockWorker) HealthCheck() error {
	if m.healthErr != nil {
		return m.healthErr
	}
	return nil
}
func (m *MockWorker) ID() string { return m.id }

func TestNewWorkerDistributor(t *testing.T) {
	t.Parallel()

	wd := NewWorkerDistributor()
	assert.NotNil(t, wd.pyManager, "Python manager should be initialized")
	assert.NotNil(t, wd.goManager, "Go manager should be initialized")
}

func TestDistributeWorkers(t *testing.T) {
	t.Parallel()

	t.Run("SuccessfulDistribution", func(t *testing.T) {
		wd := NewWorkerDistributor()
		cancel, err := wd.DistributeWorkers(4, "local", false)
		require.NoError(t, err, "Should distribute workers successfully")
		defer cancel()

		assert.Equal(t, 2, len(wd.pyManager.workers), "Should create 2 Python workers")
		assert.Equal(t, 2, len(wd.goManager.workers), "Should create 2 Go workers")
	})

	t.Run("OddWorkerCount", func(t *testing.T) {
		wd := NewWorkerDistributor()
		cancel, err := wd.DistributeWorkers(5, "local", false)
		require.NoError(t, err)
		defer cancel()

		assert.Equal(t, 2, len(wd.pyManager.workers), "Python workers should round down")
		assert.Equal(t, 3, len(wd.goManager.workers), "Go workers should get remainder")
	})

	t.Run("InvalidWorkerCount", func(t *testing.T) {
		wd := NewWorkerDistributor()
		_, err := wd.DistributeWorkers(0, "local", false)
		assert.Error(t, err, "Should reject zero workers")
	})

	t.Run("CPUResourceLimit", func(t *testing.T) {
		wd := NewWorkerDistributor()
		_, err := wd.DistributeWorkers(9999, "local", false) // Exceeds core count
		require.NoError(t, err)
		assert.LessOrEqual(t, len(wd.pyManager.workers)+len(wd.goManager.workers), runtime.NumCPU())
	})
}

func TestWorkerStartupFailures(t *testing.T) {
	t.Parallel()

	t.Run("PythonStartupFailure", func(t *testing.T) {
		wd := NewWorkerDistributor()

		// Replace manager with mock that fails
		wd.pyManager = &WorkerManager{
			workers: []Worker{&MockWorker{startErr: errors.New("python failed")}},
		}

		_, err := wd.DistributeWorkers(2, "local", false)
		assert.ErrorContains(t, err, "python worker startup failed")
	})

	t.Run("GoStartupFailure", func(t *testing.T) {
		wd := NewWorkerDistributor()

		// Setup successful Python workers
		wd.pyManager.workers = []Worker{&MockWorker{id: "py1"}}

		// Force Go failure
		wd.goManager = &WorkerManager{
			workers: []Worker{&MockWorker{startErr: errors.New("go failed")}},
		}

		_, err := wd.DistributeWorkers(2, "redis", false)
		assert.ErrorContains(t, err, "go worker startup failed")
	})
}

func TestShutdown(t *testing.T) {
	t.Parallel()

	t.Run("GracefulShutdown", func(t *testing.T) {
		wd := NewWorkerDistributor()

		// Setup mock workers
		wd.pyManager.workers = []Worker{&MockWorker{id: "py1"}}
		wd.goManager.workers = []Worker{&MockWorker{id: "go1"}}
		wd.wg.Add(2) // Simulate running workers

		// Test shutdown in separate goroutine
		go func() {
			time.Sleep(100 * time.Millisecond)
			wd.wg.Done()
			wd.wg.Done()
		}()

		err := wd.Shutdown()
		assert.NoError(t, err, "Should shutdown gracefully")
	})

	t.Run("ShutdownTimeout", func(t *testing.T) {
		wd := NewWorkerDistributor()
		wd.wg.Add(1) // Block shutdown

		err := wd.Shutdown()
		assert.ErrorContains(t, err, "shutdown timeout")
	})
}

func TestConcurrentOperations(t *testing.T) {
	t.Parallel()

	wd := NewWorkerDistributor()
	var wg sync.WaitGroup

	// Test concurrent AddWorker calls
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wd.pyManager.AddWorker(&MockWorker{id: uuid.New().String()})
			wd.goManager.AddWorker(&MockWorker{id: uuid.New().String()})
		}()
	}

	wg.Wait()
	assert.Equal(t, 10, len(wd.pyManager.workers), "Should handle concurrent adds")
	assert.Equal(t, 10, len(wd.goManager.workers), "Should handle concurrent adds")
}
