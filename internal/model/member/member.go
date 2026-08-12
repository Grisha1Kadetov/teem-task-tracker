package member

import (
	"github.com/Grisha1Kadetov/TeemTaskTrackerService/internal/model/role"
	"github.com/google/uuid"
)

type Member struct {
	TeamID uuid.UUID
	UserID uuid.UUID
	Role   role.Role
}
