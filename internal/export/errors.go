package export

import (
	"errors"
)

var (
	ErrNotFound              = errors.New("export: not found")
	ErrInvalidState          = errors.New("export: invalid state transition")
	ErrBatchMissingStartedAt = errors.New("export: batch rejected, missing started at")
	ErrBatchEmptyKind        = errors.New("export: batch rejected, empty kind")
)
