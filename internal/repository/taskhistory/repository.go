package taskhistory

import (
	"context"
	"fmt"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/taskhistory"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/db"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var queryBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Question)

type Repository struct {
	database db.Executor
}

func New(database db.Executor) *Repository {
	return &Repository{database: database}
}

func (r *Repository) Create(ctx context.Context, history taskhistory.TaskHistory) error {
	query, args, err := queryBuilder.
		Insert("task_history").
		Columns("id", "task_id", "changed_by", "changes", "created_at").
		Values(
			sq.Expr("UUID_TO_BIN(?)", history.ID),
			sq.Expr("UUID_TO_BIN(?)", history.TaskID),
			sq.Expr("UUID_TO_BIN(?)", history.ChangedBy),
			history.Changes,
			history.CreatedAt,
		).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := r.database.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("create task history: %w", err)
	}

	return nil
}

func (r *Repository) ListByTaskID(
	ctx context.Context,
	taskID uuid.UUID,
) ([]taskhistory.TaskHistory, error) {
	query, args, err := queryBuilder.
		Select("id", "task_id", "changed_by", "changes", "created_at").
		From("task_history").
		Where("task_id = UUID_TO_BIN(?)", taskID).
		OrderBy("created_at ASC", "id ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list task history: %w", err)
	}
	defer rows.Close()

	history := make([]taskhistory.TaskHistory, 0)
	for rows.Next() {
		var value taskhistory.TaskHistory
		if err := rows.Scan(
			&value.ID,
			&value.TaskID,
			&value.ChangedBy,
			&value.Changes,
			&value.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan task history: %w", err)
		}

		history = append(history, value)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate task history: %w", err)
	}

	return history, nil
}
