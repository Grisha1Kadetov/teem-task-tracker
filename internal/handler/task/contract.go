package task

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/task"
	"github.com/google/uuid"
)

type service interface {
	CreateTask(
		ctx context.Context,
		teamID uuid.UUID,
		title,
		description string,
		status task.Status,
		assigneeID uuid.UUID,
		createdBy uuid.UUID,
	) (task.Task, error)
	ListTasks(ctx context.Context, filter task.Filter) ([]task.Task, error)
	UpdateTask(
		ctx context.Context,
		title,
		description *string,
		status *task.Status,
		assigneeID *uuid.UUID,
		taskID,
		changedBy uuid.UUID,
	) error
}
