package workers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
)

type PythonWorker struct {
	cmd    *exec.Cmd
	ctx    context.Context
	cancel context.CancelFunc
	id     string
	config PythonWorkerConfig
	mu     sync.Mutex
}

type PythonWorkerConfig struct {
	EnvPath      string
	FilePath     string
	Mode         string
	IsProduction bool
}

func NewPythonWorker(parentCtx context.Context, id string, config PythonWorkerConfig) *PythonWorker {
	ctx, cancel := context.WithCancel(parentCtx)
	return &PythonWorker{
		id:     id,
		ctx:    ctx,
		cancel: cancel,
		config: config,
	}
}

func (pw *PythonWorker) Start() error {
	pw.cmd = exec.CommandContext(pw.ctx, path.Join(pw.config.EnvPath, "python3"), pw.config.FilePath, "--mode="+pw.config.Mode)

	pw.cmd.Env = append(os.Environ(),
		"PYTHONUNBUFFERED=1",
		"WORKER_ID="+pw.id,
	)

	if !pw.config.IsProduction {
		pw.cmd.Stdout = os.Stdout
		pw.cmd.Stderr = os.Stderr
	}

	return pw.cmd.Start()
}

func (pw *PythonWorker) Stop() error {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	pw.cancel()

	if pw.cmd != nil && pw.cmd.Process != nil {
		return pw.cmd.Process.Signal(os.Interrupt)
	}
	return nil
}

func (pw *PythonWorker) HealthCheck() error {
	if pw.cmd == nil || pw.cmd.Process == nil {
		return fmt.Errorf("worker not running")
	}

	if pw.cmd.ProcessState != nil && pw.cmd.ProcessState.Exited() {
		return fmt.Errorf("worker has exited")
	}
	return nil
}

func (pw *PythonWorker) ID() string {
	return pw.id
}
