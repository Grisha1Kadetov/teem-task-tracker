package member

import (
	"github.com/Grisha1Kadetov/TeamTaskTrackerService/internal/model/role"
	"github.com/google/uuid"
)

type listResponse struct {
	TeamID uuid.UUID `json:"teamID"`
	Role   role.Role `json:"role"`
}

type inviteRequest struct {
	UserID uuid.UUID `json:"userId"`
	Role   role.Role `json:"role"`
}

type inviteResponse struct {
	TeamID uuid.UUID `json:"teamID"`
	UserID uuid.UUID `json:"userId"`
	Role   role.Role `json:"role"`
}
