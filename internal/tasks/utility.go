package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/Yulian302/qugopy/config"
	"github.com/Yulian302/qugopy/internal/queue"
	"github.com/Yulian302/qugopy/logging"
	"github.com/Yulian302/qugopy/models"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

func validateTask(task models.Task) error {
	if task.Type == "" {
		return fmt.Errorf("task type cannot be empty")
	}

	if task.Priority <= 0 {
		return fmt.Errorf("priority must be positive")
	}

	if string(task.Payload) == "" {
		return fmt.Errorf("someField is required")
	}

	return nil
}

func EnqueueTask(task models.Task, rdb *redis.Client) error {
	err := validateTask(task)
	if err != nil {
		return fmt.Errorf("invalid task: %w", err)
	}
	internalTask := &models.IntTask{
		Task: task,
		ID:   uuid.New().String(),
	}
	userTaskJson, err := json.Marshal(internalTask)
	if err != nil {
		logging.DebugLog(fmt.Sprintf("Failed to marshal task: %v", err))
		return fmt.Errorf("marshal error: %w", err)
	}

	// push to redis
	if config.AppConfig.MODE == "redis" {
		return rdb.ZAdd("task_queue", redis.Z{
			Score:  float64(task.Priority),
			Member: userTaskJson,
		}).Err()
	} else {
		// enqueue locally
		queue.DefaultLocalQueue.Lock.Lock()
		defer queue.DefaultLocalQueue.Lock.Unlock()
		queue.DefaultLocalQueue.PQ.Push(*internalTask)
		return nil
	}
}
