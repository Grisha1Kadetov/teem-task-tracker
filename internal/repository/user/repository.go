package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	model "github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/user"
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

func (r *Repository) Create(ctx context.Context, user model.User) error {
	query, args, err := queryBuilder.
		Insert("users").
		Columns("id", "email", "password_hash", "name", "created_at").
		Values(
			sq.Expr("UUID_TO_BIN(?)", user.ID),
			user.Email,
			user.PasswordHash,
			user.Name,
			user.CreatedAt,
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.database.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *Repository) EmailExists(ctx context.Context, email string) (bool, error) {
	query, args, err := queryBuilder.
		Select("COUNT(*) > 0").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return false, err
	}

	var exists bool
	if err := r.database.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		return false, fmt.Errorf("check user email existence: %w", err)
	}

	return exists, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	query, args, err := queryBuilder.
		Select("id", "email", "password_hash", "name", "created_at").
		From("users").
		Where(sq.Eq{"email": email}).
		Limit(1).
		ToSql()
	if err != nil {
		return model.User{}, err
	}

	var user model.User
	err = r.database.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, model.ErrUserNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}
