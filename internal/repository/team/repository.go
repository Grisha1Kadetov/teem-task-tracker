package team

import (
	"context"
	"fmt"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/team"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/db"
	sq "github.com/Masterminds/squirrel"
)

var queryBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Question)

type Repository struct {
	database db.Database
}

func New(database db.Database) *Repository {
	return &Repository{database: database}
}

func (r *Repository) Create(ctx context.Context, team team.Team) error {
	return r.database.WithinTransaction(ctx, func(ctx context.Context) error {
		query, args, err := queryBuilder.
			Insert("teams").
			Columns("id", "name", "created_by", "created_at").
			Values(
				sq.Expr("UUID_TO_BIN(?)", team.ID),
				team.Name,
				sq.Expr("UUID_TO_BIN(?)", team.CreatedBy),
				team.CreatedAt,
			).
			ToSql()
		if err != nil {
			return err
		}

		if _, err := r.database.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("create team: %w", err)
		}

		query, args, err = queryBuilder.
			Insert("team_members").
			Columns("team_id", "user_id", "role").
			Values(
				sq.Expr("UUID_TO_BIN(?)", team.ID),
				sq.Expr("UUID_TO_BIN(?)", team.CreatedBy),
				role.Owner,
			).
			ToSql()
		if err != nil {
			return err
		}

		if _, err := r.database.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("create team owner: %w", err)
		}

		return nil
	})
}
