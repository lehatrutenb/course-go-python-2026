package errors

import "fmt"

// AdditionalMessageError ошибка, которая
type AdditionalMessageError struct {
	err    error
	format string
	args   []any
}

// NewAdditionalMessageError создает ошибку
func NewAdditionalMessageError(err error, format string, args ...any) error {
	return &AdditionalMessageError{err: err, format: format, args: args}
}

// Great?
func (e *AdditionalMessageError) Error() string {
	var fomatted string = fmt.Sprintf(e.format, e.args...)
	var errMsg string
	if e.err != nil {
		errMsg = e.err.Error()
	}
	if fomatted == "" || errMsg == "" {
		return fomatted + errMsg
	}
	return fomatted + ": " + errMsg
}

func (e *AdditionalMessageError) Unwrap() error { return e.err }

var _ error = (*AdditionalMessageError)(nil)
