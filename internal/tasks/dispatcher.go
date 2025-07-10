package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Yulian302/qugopy/internal/tasks/handlers"
	"github.com/Yulian302/qugopy/models"
)

// for tasks execution by Go workers
func DispatchTask(ctx context.Context, intTask models.IntTask) error {
	task := intTask.Task
	switch task.Type {
	case "download_file":
		var payload handlers.DownloadFilePayload
		if err := json.Unmarshal([]byte(intTask.Task.Payload), &payload); err != nil {
			return fmt.Errorf("invalid payload for download_file: %w", err)
		}

		return handlers.DownloadFile(ctx, payload.Url, payload.Filename)
	default:
		return fmt.Errorf("unknown task type: %s", task.Type)
	}
}
