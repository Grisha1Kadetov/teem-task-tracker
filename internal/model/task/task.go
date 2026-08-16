package task

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	Todo       Status = "todo"
	InProgress Status = "in_progress"
	Done       Status = "done"
)

func (s Status) IsValid() bool {
	switch s {
	case Todo, InProgress, Done:
		return true
	default:
		return false
	}
}

func ParseStatus(value string) (Status, error) {
	status := Status(value)
	if !status.IsValid() {
		return "", fmt.Errorf("unknown task status %q", value)
	}

	return status, nil
}

type Task struct {
	ID          uuid.UUID
	TeamID      uuid.UUID
	Title       string
	Description string
	Status      Status
	CreatedBy   uuid.UUID
	AssigneeID  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ClosedAt    *time.Time
	Version     uint64
}
