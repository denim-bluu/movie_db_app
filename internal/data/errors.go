package data

import "errors"

var (
	ErrInvalidRuntimeFormat = errors.New("invalid runtime format")
	ErrEditConflict         = errors.New("edit conflict")
	ErrRecordNotFound       = errors.New("record not found")
)
