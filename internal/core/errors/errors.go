package core_errors

import "errors"

var (
	ErrConflict        = errors.New("conflict")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNotFound        = errors.New("not found")
)
