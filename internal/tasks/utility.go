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

type QueueType string

const (
	PyQueue QueueType = "python_queue"
	GoQueue QueueType = "go_queue"
)

var (
	TaskQueueDict map[models.TaskType]QueueType = map[models.TaskType]QueueType{
		models.SendEmail:    GoQueue,
		models.DownloadFile: GoQueue,
		models.ProcessImage: PyQueue,
	}
)

func GetQueueType(taskType string) (QueueType, error) {
	tt := models.TaskType(taskType)
	if !tt.IsValid() {
		return "", fmt.Errorf("invalid task type: %s", taskType)
	}
	queueType, exists := TaskQueueDict[tt]
	if !exists {
		return GoQueue, nil // default queue for unknown task types
	}
	return queueType, nil
}

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
		queueType, err := GetQueueType(task.Type)
		if err != nil {
			return fmt.Errorf("wrong task type")
		}
		return rdb.ZAdd(string(queueType), redis.Z{
			Score:  float64(task.Priority),
			Member: userTaskJson,
		}).Err()
	} else {
		// enqueue locally
		queueType, err := GetQueueType(task.Type)
		if err != nil {
			return fmt.Errorf("invalid task type: %w", err)
		}
		if queueType == PyQueue {
			queue.PythonLocalQueue.Lock.Lock()
			defer queue.PythonLocalQueue.Lock.Unlock()
			queue.PythonLocalQueue.PQ.Push(*internalTask)
		} else {
			queue.GoLocalQueue.Lock.Lock()
			defer queue.GoLocalQueue.Lock.Unlock()
			queue.GoLocalQueue.PQ.Push(*internalTask)
		}
		return nil
	}
}
