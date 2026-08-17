package taskhistory

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/taskhistory"
	"github.com/google/uuid"
)

type historyRepo interface {
	Create(ctx context.Context, history taskhistory.TaskHistory) error
	ListByTaskID(ctx context.Context, taskID uuid.UUID) ([]taskhistory.TaskHistory, error)
}

type findTaskTeamService interface {
	FindTaskTeam(ctx context.Context, taskID uuid.UUID) (uuid.UUID, bool, error)
}
