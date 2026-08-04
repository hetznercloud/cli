package util

import "errors"

type ExitError struct {
	code   int
	err    error
	silent bool
}

func NewExitError(code int, err error, silent bool) *ExitError {
	if code <= 0 {
		code = 1
	}
	return &ExitError{code: code, err: err, silent: silent}
}

func (e *ExitError) Error() string {
	return e.err.Error()
}

func (e *ExitError) Unwrap() error {
	return e.err
}

func (e *ExitError) ExitCode() int {
	return e.code
}

func (e *ExitError) Silent() bool {
	return e.silent
}

func ExitCode(err error) int {
	var exitCoder interface{ ExitCode() int }
	if errors.As(err, &exitCoder) && exitCoder.ExitCode() > 0 {
		return exitCoder.ExitCode()
	}
	return 1
}

func IsSilent(err error) bool {
	var silent interface{ Silent() bool }
	return errors.As(err, &silent) && silent.Silent()
}
