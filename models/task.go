package models

import (
	"encoding/json"
	"time"

	_ "github.com/go-playground/validator"
)

const (
	SendEmail    TaskType = "send_email"
	DownloadFile TaskType = "download_file"
	ProcessImage TaskType = "process_image"
)

type Task struct {
	// Type categorizes the task (e.g., "email", "notification").
	Type string `form:"type" json:"type" binding:"required,oneof=download_file send_email process_image"`

	// Payload contains task-specific data in string format.
	Payload json.RawMessage `form:"payload" json:"payload" binding:"required"`

	// Priority determines execution order (higher values = higher priority).
	// Range: 0-65535 (uint16)
	Priority uint16 `form:"priority" json:"priority" binding:"required,min=1,max=1000"`

	// Deadline determines the deadline for a task. Task cannot be executed after the deadline. Optional field.
	Deadline *time.Time `form:"deadline" json:"deadline,omitempty"`

	// Recurring sets if a task must recur occasionally. Optional field.
	Recurring *bool `form:"recurring" json:"recurring,omitempty"`
}

type TaskType string

func (tt TaskType) IsValid() bool {
	switch tt {
	case SendEmail, DownloadFile, ProcessImage:
		return true
	default:
		return false
	}
}
