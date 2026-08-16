package taskhistory

import (
	"context"
	"fmt"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/taskhistory"
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
