package task

import "errors"

var (
	ErrInvalidStatus            = errors.New("invalid task status")
	ErrInvalidAssigneeID        = errors.New("invalid assignee ID")
	ErrAssigneeNotTeamMember    = errors.New("assignee is not a team member")
	ErrEditorNotTeamMember      = errors.New("editor is not a team member")
	ErrTaskNotFound             = errors.New("task not found")
	ErrNoChanges                = errors.New("task has no changes")
	ErrNoPermissionToUpdateTask = errors.New("user does not have permission to update task")
	ErrInsufficientPermissions  = errors.New("user does not have sufficient permissions")
)
