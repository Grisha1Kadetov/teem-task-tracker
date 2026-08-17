package task

import (
	"context"
	"time"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/task"
	"github.com/google/uuid"
)

type taskRepo interface {
	Create(ctx context.Context, task task.Task) error
	GetByID(ctx context.Context, taskID uuid.UUID) (task.Task, bool, error)
	GetByIDForUpdate(ctx context.Context, taskID uuid.UUID) (task.Task, bool, error)
	Update(ctx context.Context, taskID uuid.UUID, update TaskPatch) (bool, error)
	ListWithFilter(ctx context.Context, filter task.Filter) ([]task.Task, error)
	FindTaskTeam(ctx context.Context, taskID uuid.UUID) (uuid.UUID, bool, error)
}

type memberService interface {
	Role(ctx context.Context, userID, teamID uuid.UUID) (role.Role, bool, error)
}

type historyService interface {
	Record(
		ctx context.Context,
		before *task.Task,
		after task.Task,
		changedBy uuid.UUID,
		changedAt time.Time,
	) error
}

type transactor interface {
	WithinTransaction(ctx context.Context, fn func(context.Context) error) error
}

type UpdateTaskParams struct {
	Title       *string
	Description *string
	Status      *task.Status
	AssigneeID  *uuid.UUID
}

type TaskPatch struct {
	Title       *string
	Description *string
	Status      *task.Status
	AssigneeID  *uuid.UUID
	UpdatedAt   *time.Time
	ClosedAt    **time.Time
	Version     *uint64
}
