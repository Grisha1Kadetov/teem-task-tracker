package taskhistory

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/taskhistory"
	"github.com/google/uuid"
)

type service interface {
	ListTaskHistory(ctx context.Context, taskID uuid.UUID) ([]taskhistory.TaskHistory, error)
}
