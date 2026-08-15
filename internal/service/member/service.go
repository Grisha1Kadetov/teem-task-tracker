package member

import (
	"context"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/member"
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
	"github.com/google/uuid"
)

type MemberService struct {
	memberRepo memberRepo
}

func New(memberRepo memberRepo) *MemberService {
	return &MemberService{memberRepo: memberRepo}
}

func (s *MemberService) List(ctx context.Context, userID uuid.UUID) ([]member.Member, error) {
	return s.memberRepo.ListByUserID(ctx, userID)
}

func (s *MemberService) Invite(
	ctx context.Context,
	teamID,
	userID uuid.UUID,
	memberRole role.Role,
) (member.Member, error) {
	if memberRole != role.Admin && memberRole != role.Member {
		return member.Member{}, member.ErrRoleMismatch
	}

	teamMember := member.Member{
		TeamID: teamID,
		UserID: userID,
		Role:   memberRole,
	}
	if err := s.memberRepo.Create(ctx, teamMember); err != nil {
		return member.Member{}, err
	}

	return teamMember, nil
}

func (s *MemberService) Role(
	ctx context.Context,
	userID,
	teamID uuid.UUID,
) (role.Role, bool, error) {
	return s.memberRepo.Role(ctx, userID, teamID)
}
