package team

import (
	"context"
	"time"

	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/team"
	"github.com/google/uuid"
)

type TeamService struct {
	teamRepo teamRepo
}

func New(teamRepo teamRepo) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
	}
}

func (s *TeamService) Create(ctx context.Context, name string, createdBy uuid.UUID) (team.Team, error) {
	createdTeam := team.Team{
		ID:        uuid.New(),
		Name:      name,
		CreatedBy: createdBy,
		CreatedAt: time.Now().UTC(),
	}

	err := s.teamRepo.Create(ctx, createdTeam)
	if err != nil {
		return team.Team{}, err
	}

	return createdTeam, nil
}
