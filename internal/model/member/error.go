package member

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrRoleMismatch  = errors.New("role mismatch")
	ErrAlreadyExists = errors.New("team member already exists")
)
