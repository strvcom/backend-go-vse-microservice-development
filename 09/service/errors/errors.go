package errors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserDoesntExists  = errors.New("user does not exist")
)
