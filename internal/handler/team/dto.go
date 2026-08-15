package team

import (
	"time"

	"github.com/google/uuid"
)

type createRequest struct {
	Name string `json:"name"`
}

type createResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}
