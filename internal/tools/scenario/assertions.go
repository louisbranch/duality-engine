package scenario

import (
	"errors"
	"fmt"
	"log"
)

// AssertionMode controls whether validation failures error or only log.
type AssertionMode int

const (
	// AssertionStrict returns errors on failed assertions.
	AssertionStrict AssertionMode = iota
	// AssertionLogOnly logs failed assertions and continues.
	AssertionLogOnly
)

// Assertions handles scenario expectations with optional strictness.
type Assertions struct {
	Mode   AssertionMode
	Logger *log.Logger
}

func (a Assertions) Assertf(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	if a.Mode == AssertionLogOnly {
		a.logf("assertion skipped: %s", message)
		return nil
	}
	return errors.New(message)
}

func (a Assertions) Failf(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	return errors.New(message)
}

func (a Assertions) logf(format string, args ...any) {
	if a.Logger != nil {
		a.Logger.Printf(format, args...)
	}
}
