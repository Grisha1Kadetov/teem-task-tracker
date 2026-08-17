package task

import (
	"time"

	model "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/task"
	"github.com/google/uuid"
)

type TaskPatch struct {
	Title       *string
	Description *string
	Status      *model.Status
	AssigneeID  *uuid.UUID
	UpdatedAt   *time.Time
	ClosedAt    **time.Time
	Version     *uint64
}
