package postgres

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrProfileNotFound = errors.New("profile not found")
	ErrAccountNotFound = errors.New("account not found")
)