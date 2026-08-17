package taskhistory

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type historyResponse struct {
	ID        uuid.UUID       `json:"id"`
	TaskID    uuid.UUID       `json:"task_id"`
	ChangedBy uuid.UUID       `json:"changed_by"`
	Changes   json.RawMessage `json:"changes"`
	CreatedAt time.Time       `json:"created_at"`
}
