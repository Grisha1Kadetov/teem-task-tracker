package task

import (
	"context"
	"fmt"
	"time"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/task"
	"github.com/google/uuid"
)

const defaultLimit uint64 = 20

type TaskService struct {
	taskRepo       taskRepo
	memberService  memberService
	historyService historyService
	transactor     transactor
}

func New(
	taskRepo taskRepo,
	memberService memberService,
	historyService historyService,
	transactor transactor,
) *TaskService {
	return &TaskService{
		taskRepo:       taskRepo,
		memberService:  memberService,
		historyService: historyService,
		transactor:     transactor,
	}
}

func (s *TaskService) CreateTask(
	ctx context.Context,
	teamID uuid.UUID,
	title,
	description string,
	status task.Status,
	assigneeID uuid.UUID,
	createdBy uuid.UUID,
) (task.Task, error) {
	if !status.IsValid() {
		return task.Task{}, task.ErrInvalidStatus
	}
	if assigneeID == uuid.Nil {
		return task.Task{}, task.ErrInvalidAssigneeID
	}

	_, found, err := s.memberService.Role(ctx, assigneeID, teamID)
	if err != nil {
		return task.Task{}, err
	}
	if !found {
		return task.Task{}, task.ErrAssigneeNotTeamMember
	}

	now := time.Now().UTC()
	createdTask := task.Task{
		ID:          uuid.New(),
		TeamID:      teamID,
		Title:       title,
		Description: description,
		Status:      status,
		CreatedBy:   createdBy,
		AssigneeID:  assigneeID,
		CreatedAt:   now,
		UpdatedAt:   now,
		Version:     1,
	}
	if status == task.Done {
		createdTask.ClosedAt = &now
	}

	if err := s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := s.taskRepo.Create(ctx, createdTask); err != nil {
			return err
		}

		return s.historyService.Record(ctx, nil, createdTask, createdBy, now)
	}); err != nil {
		return task.Task{}, err
	}

	return createdTask, nil
}

func (s *TaskService) ListTasks(ctx context.Context, filter task.Filter) ([]task.Task, error) {
	if filter.Status != nil && !filter.Status.IsValid() {
		return nil, task.ErrInvalidStatus
	}
	if filter.Limit == 0 {
		filter.Limit = defaultLimit
	}

	return s.taskRepo.ListWithFilter(ctx, filter)
}

func (s *TaskService) FindTaskTeam(ctx context.Context, taskID uuid.UUID) (uuid.UUID, bool, error) {
	return s.taskRepo.FindTaskTeam(ctx, taskID)
}

func (s *TaskService) GetTask(ctx context.Context, taskID uuid.UUID) (task.Task, error) {
	value, found, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return task.Task{}, err
	}
	if !found {
		return task.Task{}, task.ErrTaskNotFound
	}

	return value, nil
}

func (s *TaskService) UpdateTask(
	ctx context.Context,
	update UpdateTaskParams,
	taskID uuid.UUID,
	changedBy uuid.UUID,
) error {
	if changedBy == uuid.Nil {
		return fmt.Errorf("changedBy cannot be nil")
	}

	return s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		taskForUpdate, err := s.getTaskForUpdate(ctx, taskID)
		if err != nil {
			return err
		}
		if err := s.authorizeTaskUpdate(ctx, taskForUpdate, update, changedBy); err != nil {
			return err
		}
		if err := s.validateTaskUpdate(ctx, taskForUpdate, update); err != nil {
			return err
		}

		now := time.Now().UTC()
		updatedTask, patch, err := buildTaskUpdate(taskForUpdate, update, now)
		if err != nil {
			return err
		}

		updated, err := s.taskRepo.Update(ctx, taskForUpdate.ID, patch)
		if err != nil {
			return err
		}
		if !updated {
			return task.ErrTaskNotFound
		}

		return s.historyService.Record(ctx, &taskForUpdate, updatedTask, changedBy, now)
	})
}

func (s *TaskService) getTaskForUpdate(ctx context.Context, taskID uuid.UUID) (task.Task, error) {
	value, found, err := s.taskRepo.GetByIDForUpdate(ctx, taskID)
	if err != nil {
		return task.Task{}, err
	}
	if !found {
		return task.Task{}, task.ErrTaskNotFound
	}

	return value, nil
}

func (s *TaskService) authorizeTaskUpdate(
	ctx context.Context,
	taskForUpdate task.Task,
	update UpdateTaskParams,
	changedBy uuid.UUID,
) error {
	editorRole, found, err := s.memberService.Role(ctx, changedBy, taskForUpdate.TeamID)
	if err != nil {
		return err
	}
	if !found {
		return task.ErrEditorNotTeamMember
	}

	canUpdateAll := taskForUpdate.CreatedBy == changedBy ||
		editorRole == role.Owner ||
		editorRole == role.Admin

	if !canUpdateAll && taskForUpdate.AssigneeID != changedBy {
		return task.ErrNoPermissionToUpdateTask
	}
	if !canUpdateAll && (update.Title != nil || update.Description != nil || update.AssigneeID != nil) {
		return task.ErrInsufficientPermissions
	}

	return nil
}

func (s *TaskService) validateTaskUpdate(
	ctx context.Context,
	taskForUpdate task.Task,
	update UpdateTaskParams,
) error {
	if update.Status != nil && !update.Status.IsValid() {
		return task.ErrInvalidStatus
	}
	if update.AssigneeID != nil {
		if *update.AssigneeID == uuid.Nil {
			return task.ErrInvalidAssigneeID
		}
		if *update.AssigneeID != taskForUpdate.AssigneeID {
			_, found, err := s.memberService.Role(ctx, *update.AssigneeID, taskForUpdate.TeamID)
			if err != nil {
				return err
			}
			if !found {
				return task.ErrAssigneeNotTeamMember
			}
		}
	}

	return nil
}

func buildTaskUpdate(
	taskForUpdate task.Task,
	update UpdateTaskParams,
	now time.Time,
) (task.Task, TaskPatch, error) {
	updatedTask := taskForUpdate
	patch := TaskPatch{}
	changed := false
	if update.Title != nil && *update.Title != updatedTask.Title {
		updatedTask.Title = *update.Title
		patch.Title = update.Title
		changed = true
	}
	if update.Description != nil && *update.Description != updatedTask.Description {
		updatedTask.Description = *update.Description
		patch.Description = update.Description
		changed = true
	}
	if update.Status != nil && *update.Status != updatedTask.Status {
		updatedTask.Status = *update.Status
		patch.Status = update.Status
		changed = true
	}
	if update.AssigneeID != nil && *update.AssigneeID != updatedTask.AssigneeID {
		updatedTask.AssigneeID = *update.AssigneeID
		patch.AssigneeID = update.AssigneeID
		changed = true
	}
	if !changed {
		return task.Task{}, TaskPatch{}, task.ErrNoChanges
	}

	closedAt := updatedTask.ClosedAt
	if update.Status != nil && *update.Status != taskForUpdate.Status {
		if *update.Status == task.Done {
			closedAt = &now
		} else {
			closedAt = nil
		}
		updatedTask.ClosedAt = closedAt
		patch.ClosedAt = &closedAt
	}
	updatedTask.UpdatedAt = now
	updatedTask.Version++
	patch.UpdatedAt = &updatedTask.UpdatedAt
	patch.Version = &updatedTask.Version

	return updatedTask, patch, nil
}
