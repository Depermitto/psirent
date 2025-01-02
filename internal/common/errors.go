package common

import "errors"

var (
	ErrFileNotShared = errors.New("could not share the file")
	ErrInvalidAddr   = errors.New("invalid addr")
	ErrDuplicate     = errors.New("same peer shared the same file")
)
