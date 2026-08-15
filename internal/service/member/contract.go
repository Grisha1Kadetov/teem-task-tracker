package member

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/member"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
	"github.com/google/uuid"
)

type memberRepo interface {
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]member.Member, error)
	Create(ctx context.Context, member member.Member) error
	Role(ctx context.Context, userID, teamID uuid.UUID) (role.Role, bool, error)
}
