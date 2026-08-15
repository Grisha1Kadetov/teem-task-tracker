package teamaccess

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
	"github.com/google/uuid"
)

type service interface {
	Role(ctx context.Context, userID, teamID uuid.UUID) (role.Role, bool, error)
}

type FindTeamFunc func(ctx context.Context, resourceID uuid.UUID) (uuid.UUID, bool, error)
