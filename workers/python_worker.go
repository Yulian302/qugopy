package workers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
)

type PythonWorker struct {
	cmd        *exec.Cmd
	ctx        context.Context
	cancel     context.CancelFunc
	id         string
	config     PythonWorkerConfig
	mu         sync.Mutex
	waitDoneCh chan struct{}
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
		id:         id,
		ctx:        ctx,
		cancel:     cancel,
		config:     config,
		waitDoneCh: make(chan struct{}),
	}
}

func (pw *PythonWorker) Start() error {
	pw.cmd = exec.CommandContext(pw.ctx, path.Join(pw.config.EnvPath, "python3"), pw.config.FilePath)

	pw.cmd.Env = append(os.Environ(),
		"IS_PRODUCTION="+strconv.FormatBool(pw.config.IsProduction),
		"MODE="+pw.config.Mode,
		"PYTHONUNBUFFERED=1",
		"WORKER_ID="+pw.id,
	)

	if !pw.config.IsProduction {
		pw.cmd.Stdout = os.Stdout
		pw.cmd.Stderr = os.Stderr
	}

	err := pw.cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := pw.cmd.Wait()
		if err != nil {
			fmt.Printf("Python worker %s exited with error: %v\n", pw.id, err)
		} else {
			fmt.Printf("Python worker %s exited normally\n", pw.id)
		}
		close(pw.waitDoneCh)
	}()

	return nil
}

func (pw *PythonWorker) Stop() error {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	pw.cancel()

	if pw.cmd != nil && pw.cmd.Process != nil {
		if err := pw.cmd.Process.Signal(os.Interrupt); err != nil {
			fmt.Printf("Error signaling process %s: %v\n", pw.id, err)
		}
		<-pw.waitDoneCh
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
