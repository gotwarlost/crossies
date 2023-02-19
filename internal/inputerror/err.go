package inputerror

import (
	"github.com/pkg/errors"
)

type inError interface {
	isInputError()
}

type inputError struct {
	msg string
}

func (i inputError) Error() string {
	return i.msg
}

func (i inputError) isInputError() {}

func (i inputError) Is(err error) bool {
	_, ok := err.(inError)
	return ok
}

// New returns an input error with the supplied message.
func New(msg string) error {
	return inputError{msg: msg}
}

// IsInputError returns true if the supplied error is an input error.
func IsInputError(err error) bool {
	return errors.Is(err, inputError{})
}
