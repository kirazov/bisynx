package hashtree

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyTree       = errors.New("hashtree: tree cannot be empty")
	ErrAlreadyCompiled = errors.New("hashtree: already compiled")
)

type HashSizeError struct {
	expect, got int
}

func (e *HashSizeError) Error() string {
	return fmt.Sprintf(
		"hashtree: hash size mismatch: expected %d, got %d", e.expect, e.got,
	)
}
