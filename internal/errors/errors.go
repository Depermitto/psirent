package errors

import "errors"

var (
	ErrInvalidAddr     = errors.New("invalid addr")
	ErrUnknownResponse = errors.New("unknown response")
	ErrLostConnection  = errors.New("lost connection unexpectedly while sharing network resources")

	ErrShareFileNotShared = errors.New("could not share the file")
	ErrShareDuplicate     = errors.New("peer shared the same file")

	ErrLsEmpty = errors.New("no files available for download")
)
