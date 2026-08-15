package teamaccess

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
	"github.com/google/uuid"
)

type Access struct {
	TeamID uuid.UUID
	Role   role.Role
}

type accessContextKey struct{}

func withAccess(ctx context.Context, access Access) context.Context {
	return context.WithValue(ctx, accessContextKey{}, access)
}

func FromContext(ctx context.Context) (Access, bool) {
	access, ok := ctx.Value(accessContextKey{}).(Access)
	return access, ok
}
