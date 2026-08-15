package team

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/team"
)

type teamRepo interface {
	Create(ctx context.Context, team team.Team) error
}
