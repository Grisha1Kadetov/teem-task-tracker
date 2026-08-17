package task

import (
	"context"

	model "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/task"
	taskService "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/service/task"
	"github.com/google/uuid"
)

type Service struct {
	*taskService.TaskService
}

func NewService(service *taskService.TaskService) *Service {
	return &Service{TaskService: service}
}

func (s *Service) UpdateTask(
	ctx context.Context,
	title,
	description *string,
	status *model.Status,
	assigneeID *uuid.UUID,
	taskID,
	changedBy uuid.UUID,
) error {
	return s.TaskService.UpdateTask(ctx, taskService.UpdateTaskParams{
		Title:       title,
		Description: description,
		Status:      status,
		AssigneeID:  assigneeID,
	}, taskID, changedBy)
}
