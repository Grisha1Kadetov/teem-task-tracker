package member

import (
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
	"github.com/google/uuid"
)

type Member struct {
	TeamID uuid.UUID
	UserID uuid.UUID
	Role   role.Role
}
