package auth

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/user"
)

type service interface {
	Register(ctx context.Context, email, password, name string) (user.User, string, error)
	Login(ctx context.Context, email, password string) (string, error)
}
