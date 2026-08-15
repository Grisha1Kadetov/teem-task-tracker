package team

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/team"
	"github.com/google/uuid"
)

type service interface {
	Create(ctx context.Context, name string, createdBy uuid.UUID) (team.Team, error)
}
