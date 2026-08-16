package taskhistory

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/task"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/taskhistory"
	"github.com/google/uuid"
)

type TaskHistoryService struct {
	historyRepo historyRepo
}

type changes struct {
	Before map[string]any `json:"before"`
	After  map[string]any `json:"after"`
}

func New(historyRepo historyRepo) *TaskHistoryService {
	return &TaskHistoryService{historyRepo: historyRepo}
}

func (s *TaskHistoryService) Record(
	ctx context.Context,
	before *task.Task,
	after task.Task,
	changedBy uuid.UUID,
	changedAt time.Time,
) error {
	beforeValues := snapshot(before)
	afterValues := snapshot(&after)

	for field, afterValue := range afterValues {
		beforeValue, exists := beforeValues[field]
		if !exists {
			if afterValue == nil {
				delete(afterValues, field)
			}
			continue
		}
		if reflect.DeepEqual(beforeValue, afterValue) {
			delete(beforeValues, field)
			delete(afterValues, field)
		}
	}
	if len(afterValues) == 0 {
		return taskhistory.ErrNoChanges
	}

	encodedChanges, err := json.Marshal(changes{
		Before: beforeValues,
		After:  afterValues,
	})
	if err != nil {
		return err
	}

	return s.historyRepo.Create(ctx, taskhistory.TaskHistory{
		ID:        uuid.New(),
		TaskID:    after.ID,
		ChangedBy: changedBy,
		Changes:   string(encodedChanges),
		CreatedAt: changedAt.UTC(),
	})
}

func snapshot(value *task.Task) map[string]any {
	if value == nil {
		return make(map[string]any)
	}

	values := map[string]any{
		"team_id":     value.TeamID.String(),
		"title":       value.Title,
		"description": value.Description,
		"status":      string(value.Status),
		"assignee_id": value.AssigneeID.String(),
		"closed_at":   nil,
	}
	if value.ClosedAt != nil {
		values["closed_at"] = value.ClosedAt.UTC().Format(time.RFC3339Nano)
	}

	return values
}
