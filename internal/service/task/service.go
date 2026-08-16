package task

import (
	"context"
	"time"

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
