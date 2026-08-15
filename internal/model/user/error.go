package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrIncorrectPassword  = errors.New("incorrect user password")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
