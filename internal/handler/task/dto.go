package task

import (
	"time"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/task"
	"github.com/google/uuid"
)

type createRequest struct {
	TeamID      uuid.UUID   `json:"team_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Status      task.Status `json:"status"`
	AssigneeID  uuid.UUID   `json:"assignee_id"`
}

type taskResponse struct {
	ID          uuid.UUID   `json:"id"`
	TeamID      uuid.UUID   `json:"team_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Status      task.Status `json:"status"`
	CreatedBy   uuid.UUID   `json:"created_by"`
	AssigneeID  uuid.UUID   `json:"assignee_id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	ClosedAt    *time.Time  `json:"closed_at"`
	Version     uint64      `json:"version"`
}
