package inputerror

// InputError signals an error in inputs.
type InputError interface {
	IsInputError() bool
}

type inputError struct {
	msg string
}

func (i inputError) Error() string {
	return i.msg
}

func (i inputError) IsInputError() bool {
	return true
}

// New returns an input error with the supplied message.
func New(msg string) error {
	return inputError{msg: msg}
}
