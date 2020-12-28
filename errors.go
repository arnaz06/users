package users

import (
	"errors"
	"fmt"
)

var (

	// ErrNotFound is thrown if any requested object is doesn't exists.
	ErrNotFound = errors.New("Your requested object does not exists")
)

// ConstraintError represents a custom error for a contstraint things.
type ConstraintError string

func (e ConstraintError) Error() string {
	return string(e)
}

// ConstraintErrorf constructs ConstraintError with formatted message.
func ConstraintErrorf(format string, a ...interface{}) ConstraintError {
	return ConstraintError(fmt.Sprintf(format, a...))
}

// UnauthorizedError represents a custom error for an error related with authentication things.
type UnauthorizedError string

func (e UnauthorizedError) Error() string {
	return string(e)
}

// UnauthorizedErrorf constructs UnauthorizedError with formatted message.
func UnauthorizedErrorf(format string, a ...interface{}) UnauthorizedError {
	return UnauthorizedError(fmt.Sprintf(format, a...))
}
