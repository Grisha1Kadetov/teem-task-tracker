package taskhistory

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/taskhistory"
)

type historyRepo interface {
	Create(ctx context.Context, history taskhistory.TaskHistory) error
}
