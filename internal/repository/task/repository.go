package task

import (
	"context"
	"database/sql"
	"fmt"

	model "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/task"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/db"
	sq "github.com/Masterminds/squirrel"
)

var queryBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Question)

type Repository struct {
	database db.Executor
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
		var value model.Task
		var rawStatus string
		var closedAt sql.NullTime
		if err := rows.Scan(
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
			return nil, fmt.Errorf("scan task: %w", err)
		}

		status, err := model.ParseStatus(rawStatus)
		if err != nil {
			return nil, fmt.Errorf("parse task status: %w", err)
		}
		value.Status = status
		if closedAt.Valid {
			value.ClosedAt = &closedAt.Time
		}

		tasks = append(tasks, value)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}
