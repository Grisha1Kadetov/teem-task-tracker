package auth

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/user"
)

type userRepo interface {
	Create(ctx context.Context, user user.User) error
	EmailExists(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (user.User, error)
}
