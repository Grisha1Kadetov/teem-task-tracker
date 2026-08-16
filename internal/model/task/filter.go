package task

import "github.com/google/uuid"

type Filter struct {
	TeamID     uuid.UUID
	Status     *Status
	AssigneeID *uuid.UUID
	Limit      uint64
	Offset     uint64
}
