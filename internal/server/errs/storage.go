package errs

import "errors"

var (
	ErrStorageInvalidGaugeName   = errors.New("invalid gauge name")
	ErrStorageInvalidCounterName = errors.New("invalid counter name")
)
