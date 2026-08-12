package taskhistory

import (
	"time"

	"github.com/google/uuid"
)

type TaskHistory struct {
	ID        uuid.UUID
	TaskID    uuid.UUID
	ChangedBy uuid.UUID
	Changes   string
	CreatedAt time.Time
}
