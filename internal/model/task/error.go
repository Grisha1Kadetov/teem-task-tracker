package task

import "errors"

var (
	ErrInvalidStatus         = errors.New("invalid task status")
	ErrInvalidAssigneeID     = errors.New("invalid assignee ID")
	ErrAssigneeNotTeamMember = errors.New("assignee is not a team member")
)
