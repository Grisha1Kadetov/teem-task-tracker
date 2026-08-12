package task

import (
	"time"

	"github.com/google/uuid"
)

type Status string

type Task struct {
	ID          uuid.UUID
	TeamID      uuid.UUID
	Title       string
	Description string
	Status      Status
	CreatedBy   uuid.UUID
	AssigneeID  *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ClosedAt    *time.Time
	Version     uint64
}
