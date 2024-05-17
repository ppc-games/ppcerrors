package ppcerrors

import "strings"

type (
	// ErrorCoder interface defines the methods that an error code must implement.
	ErrorCoder interface {
		Name() string
		Code() int
		Msg() string
	}

	// errorCode defines an error with a name, code, and message.
	errorCode struct {
		name string
		code int
		msg  string
	}
)

// NewErrorCode creates and returns a pointer to an error code instance.
func NewErrorCode(name string, code int, msg string) *errorCode {
	return &errorCode{name: name, code: code, msg: msg}
}

func (c *errorCode) Name() string {
	return c.name
}

func (c *errorCode) Code() int {
	return c.code
}

func (c *errorCode) Msg() string {
	return c.msg
}

// New creates a new error with the given messages and associates it with the error code.
// It returns an error that implements the `error` interface,
// when Config.Caller == true, pc records the function name, file, and line number of the method that called this method.
func (c *errorCode) New(messages ...string) error {
	return &withErrorCode{
		errCode: c,
		msg:     strings.Join(messages, messagesSeparator),
		pc:      getPCFromCaller(),
	}
}

// Wrap wraps the given error with additional context and returns a new error.
// If the cause is nil, it returns nil.
// The additional context is specified by the messages parameter, which is joined
// using the messagesSeparator. The function also captures the program counter (PC)
// of the caller using the getPCFromCaller function.
func (c *errorCode) Wrap(cause error, messages ...string) error {
	if cause == nil {
		return nil
	}

	return &withCause{
		error: &withErrorCode{
			errCode: c,
			msg:     strings.Join(messages, messagesSeparator),
			pc:      getPCFromCaller(),
		},
		cause: cause,
	}
}
