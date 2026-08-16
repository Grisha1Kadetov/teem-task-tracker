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
	ListWithFilter(ctx context.Context, filter task.Filter) ([]task.Task, error)
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
