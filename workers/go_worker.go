package workers

import (
	"context"
	"fmt"
	"log"
	"time"
)

type GoWorker struct {
	ctx        context.Context
	cancel     context.CancelFunc
	id         string
	workerFunc func(context.Context) error
	done       chan struct{}
}

func NewGoWorker(id string, workerFunc func(context.Context) error) *GoWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &GoWorker{
		id:         id,
		ctx:        ctx,
		cancel:     cancel,
		workerFunc: workerFunc,
		done:       make(chan struct{}),
	}
}

func (gw *GoWorker) Start() error {
	go func() {
		defer close(gw.done)
		err := gw.workerFunc(gw.ctx)
		if err != nil && err != context.Canceled {
			log.Printf("Worker %s exited with error: %v", gw.id, err)
		}
	}()
	return nil
}

func (gw *GoWorker) Stop() error {
	gw.cancel()
	select {
	case <-gw.done:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for worker to stop")
	}
}

func (gw *GoWorker) HealthCheck() error {
	select {
	case <-gw.done:
		return fmt.Errorf("worker has exited")
	default:
		return nil
	}
}

func (gw *GoWorker) ID() string {
	return gw.id
}
