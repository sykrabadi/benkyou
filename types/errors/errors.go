package errors

import "errors"

var (
	ErrDirNotFound  = errors.New("directory not found")
	ErrFileNotFound = errors.New("file not found")
	ErrLevelDoesNotExists = errors.New("level not found")
	ErrInsufficientOptions = errors.New("insufficient options")
)
