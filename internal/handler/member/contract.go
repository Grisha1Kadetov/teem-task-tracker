package member

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/member"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
	"github.com/google/uuid"
)

type service interface {
	List(ctx context.Context, userID uuid.UUID) ([]member.Member, error)
	Invite(ctx context.Context, teamID, userID uuid.UUID, role role.Role) (member.Member, error)
}
