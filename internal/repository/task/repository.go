package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	model "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/task"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/db"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var queryBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Question)

type Repository struct {
	database db.Executor
}

type scanner interface {
	Scan(dest ...any) error
}

func New(database db.Executor) *Repository {
	return &Repository{database: database}
}

func (r *Repository) Create(ctx context.Context, createdTask model.Task) error {
	query, args, err := queryBuilder.
		Insert("tasks").
		Columns(
			"id",
			"team_id",
			"title",
			"description",
			"status",
			"created_by",
			"assignee_id",
			"created_at",
			"updated_at",
			"closed_at",
			"version",
		).
		Values(
			sq.Expr("UUID_TO_BIN(?)", createdTask.ID),
			sq.Expr("UUID_TO_BIN(?)", createdTask.TeamID),
			createdTask.Title,
			createdTask.Description,
			createdTask.Status,
			sq.Expr("UUID_TO_BIN(?)", createdTask.CreatedBy),
			sq.Expr("UUID_TO_BIN(?)", createdTask.AssigneeID),
			createdTask.CreatedAt,
			createdTask.UpdatedAt,
			createdTask.ClosedAt,
			createdTask.Version,
		).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.database.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("create task: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	taskID uuid.UUID,
) (model.Task, bool, error) {
	return r.getByID(ctx, taskID, false)
}

func (r *Repository) GetByIDForUpdate(
	ctx context.Context,
	taskID uuid.UUID,
) (model.Task, bool, error) {
	return r.getByID(ctx, taskID, true)
}

func (r *Repository) getByID(
	ctx context.Context,
	taskID uuid.UUID,
	forUpdate bool,
) (model.Task, bool, error) {
	builder := queryBuilder.
		Select(
			"id",
			"team_id",
			"title",
			"description",
			"status",
			"created_by",
			"assignee_id",
			"created_at",
			"updated_at",
			"closed_at",
			"version",
		).
		From("tasks").
		Where("id = UUID_TO_BIN(?)", taskID).
		Limit(1)
	if forUpdate {
		builder = builder.Suffix("FOR UPDATE")
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return model.Task{}, false, err
	}

	value, err := scanTask(r.database.QueryRowContext(ctx, query, args...))
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, false, nil
	}
	if err != nil {
		return model.Task{}, false, fmt.Errorf("get task: %w", err)
	}

	return value, true, nil
}

func (r *Repository) Update(
	ctx context.Context,
	taskID uuid.UUID,
	update TaskPatch,
) (bool, error) {
	builder := queryBuilder.Update("tasks")
	hasUpdates := false
	if update.Title != nil {
		builder = builder.Set("title", *update.Title)
		hasUpdates = true
	}
	if update.Description != nil {
		builder = builder.Set("description", *update.Description)
		hasUpdates = true
	}
	if update.Status != nil {
		builder = builder.Set("status", *update.Status)
		hasUpdates = true
	}
	if update.ClosedAt != nil {
		builder = builder.Set("closed_at", *update.ClosedAt)
		hasUpdates = true
	}
	if update.AssigneeID != nil {
		builder = builder.Set("assignee_id", sq.Expr("UUID_TO_BIN(?)", *update.AssigneeID))
		hasUpdates = true
	}
	if update.UpdatedAt != nil {
		builder = builder.Set("updated_at", *update.UpdatedAt)
		hasUpdates = true
	}
	if update.Version != nil {
		builder = builder.Set("version", *update.Version)
		hasUpdates = true
	}
	if !hasUpdates {
		return false, nil
	}

	query, args, err := builder.
		Where("id = UUID_TO_BIN(?)", taskID).
		ToSql()
	if err != nil {
		return false, err
	}

	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return false, fmt.Errorf("update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("get updated task rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}

func (r *Repository) ListWithFilter(ctx context.Context, filter model.Filter) ([]model.Task, error) {
	builder := queryBuilder.
		Select(
			"id",
			"team_id",
			"title",
			"description",
			"status",
			"created_by",
			"assignee_id",
			"created_at",
			"updated_at",
			"closed_at",
			"version",
		).
		From("tasks").
		Where("team_id = UUID_TO_BIN(?)", filter.TeamID).
		OrderBy("created_at DESC", "id DESC").
		Limit(filter.Limit).
		Offset(filter.Offset)
	if filter.Status != nil {
		builder = builder.Where(sq.Eq{"status": *filter.Status})
	}
	if filter.AssigneeID != nil {
		builder = builder.Where("assignee_id = UUID_TO_BIN(?)", *filter.AssigneeID)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		value, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		tasks = append(tasks, value)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}

func (r *Repository) FindTaskTeam(
	ctx context.Context,
	taskID uuid.UUID,
) (uuid.UUID, bool, error) {
	query, args, err := queryBuilder.
		Select("team_id").
		From("tasks").
		Where("id = UUID_TO_BIN(?)", taskID).
		Limit(1).
		ToSql()
	if err != nil {
		return uuid.Nil, false, err
	}

	var teamID uuid.UUID
	err = r.database.QueryRowContext(ctx, query, args...).Scan(&teamID)
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, false, nil
	}
	if err != nil {
		return uuid.Nil, false, fmt.Errorf("find task team: %w", err)
	}

	return teamID, true, nil
}

func scanTask(row scanner) (model.Task, error) {
	var value model.Task
	var rawStatus string
	var closedAt sql.NullTime
	if err := row.Scan(
		&value.ID,
		&value.TeamID,
		&value.Title,
		&value.Description,
		&rawStatus,
		&value.CreatedBy,
		&value.AssigneeID,
		&value.CreatedAt,
		&value.UpdatedAt,
		&closedAt,
		&value.Version,
	); err != nil {
		return model.Task{}, err
	}

	status, err := model.ParseStatus(rawStatus)
	if err != nil {
		return model.Task{}, fmt.Errorf("parse task status: %w", err)
	}
	value.Status = status
	if closedAt.Valid {
		value.ClosedAt = &closedAt.Time
	}

	return value, nil
}
