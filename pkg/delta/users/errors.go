package users

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidID    = errors.New("invalid user ID")
	ErrUserExists   = errors.New("user already exists")
)
