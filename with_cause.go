package ppcerrors

import (
	"fmt"
	"io"
	"strings"
)

// withCause implements the error interface.
// The embedded error field is used to store the wrapped error with possible additional information when passing the error around the system,
// and the cause field is used to wrap the root cause of the error.
type withCause struct {
	error
	cause error
}

// Error prints the error message of the current error e, followed by the error message of the cause wrapped by e.
// e.g.: ErrNilUser, User information is empty <= cause's Error().
func (e *withCause) Error() string {
	var b strings.Builder
	b.WriteString(e.error.Error())
	b.WriteString(Config.ErrorChainSeparator)
	b.WriteString(e.cause.Error())
	return b.String()
}

// Cause returns the cause wrapped by the current error e.
// It implements the causer interface in the github.com/pkg/errors package.
func (e *withCause) Cause() error {
	return e.cause
}

// Unwrap returns the cause wrapped by the current error e.
// It implements the Unwrap interface in the errors standard library.
func (e *withCause) Unwrap() error {
	return e.cause
}

// Format will print the detailed error reasons of each layer of cause in the error chain when verb == %+v.
// Otherwise, it will print the error message of the current error.
func (e *withCause) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// 先打印当前错误
			_, _ = fmt.Fprintf(s, "%+v", e.error)

			// 然后打印被包装的 cause
			_, _ = fmt.Fprintf(s, "\ncause: %+v", e.cause)

			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}
