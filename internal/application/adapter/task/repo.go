package task

import (
	"context"

	taskRepository "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/repository/task"
	taskService "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/service/task"
	"github.com/google/uuid"
)

type Adapter struct {
	*taskRepository.Repository
}

func New(repository *taskRepository.Repository) *Adapter {
	return &Adapter{Repository: repository}
}

func (a *Adapter) Update(
	ctx context.Context,
	taskID uuid.UUID,
	update taskService.TaskPatch,
) (bool, error) {
	return a.Repository.Update(ctx, taskID, taskRepository.TaskPatch{
		Title:       update.Title,
		Description: update.Description,
		Status:      update.Status,
		AssigneeID:  update.AssigneeID,
		UpdatedAt:   update.UpdatedAt,
		ClosedAt:    update.ClosedAt,
		Version:     update.Version,
	})
}
