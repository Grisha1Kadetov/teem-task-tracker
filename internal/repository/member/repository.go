package member

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	model "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/member"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/db"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/pkg/mysqlerr"
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

func (r *Repository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Member, error) {
	query, args, err := queryBuilder.
		Select("team_id", "role").
		From("team_members").
		Where("user_id = UUID_TO_BIN(?)", userID).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list user team memberships: %w", err)
	}
	defer rows.Close()

	members := make([]model.Member, 0)
	for rows.Next() {
		var teamMember model.Member
		var rawRole string
		if err := rows.Scan(&teamMember.TeamID, &rawRole); err != nil {
			return nil, fmt.Errorf("scan user team membership: %w", err)
		}

		teamRole, err := role.Parse(rawRole)
		if err != nil {
			return nil, fmt.Errorf("parse user team membership role: %w", err)
		}

		teamMember.UserID = userID
		teamMember.Role = teamRole
		members = append(members, teamMember)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user team memberships: %w", err)
	}

	return members, nil
}

func (r *Repository) Create(ctx context.Context, teamMember model.Member) error {
	selectUser := queryBuilder.
		Select().
		Column("UUID_TO_BIN(?)", teamMember.TeamID).
		Column("id").
		Column("?", teamMember.Role).
		From("users").
		Where("id = UUID_TO_BIN(?)", teamMember.UserID)

	query, args, err := queryBuilder.
		Insert("team_members").
		Columns("team_id", "user_id", "role").
		Select(selectUser).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		if mysqlerr.IsDuplicateEntry(err) {
			return model.ErrAlreadyExists
		}

		return fmt.Errorf("create team member: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get created team member count: %w", err)
	}
	if rowsAffected == 0 {
		return model.ErrUserNotFound
	}

	return nil
}

func (r *Repository) Role(ctx context.Context, userID, teamID uuid.UUID) (role.Role, bool, error) {
	query, args, err := queryBuilder.
		Select("role").
		From("team_members").
		Where("user_id = UUID_TO_BIN(?)", userID).
		Where("team_id = UUID_TO_BIN(?)", teamID).
		Limit(1).
		ToSql()
	if err != nil {
		return "", false, err
	}

	var rawRole string
	err = r.database.QueryRowContext(ctx, query, args...).Scan(&rawRole)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("get team member role: %w", err)
	}

	teamRole, err := role.Parse(rawRole)
	if err != nil {
		return "", false, fmt.Errorf("parse team member role: %w", err)
	}

	return teamRole, true, nil
}
