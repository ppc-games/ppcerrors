package ppcerrors

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	// WithErrorCoder defines the ErrorCode() method, which is used to return the error code contained in an error.
	WithErrorCoder interface {
		ErrorCode() *errorCode
	}

	// withErrorCode is an error that contains an error code to distinguish it from other errors.
	// The error code is generally used to return to systems outside the current application system boundary (e.g., clients),
	// because these external systems cannot directly get the error, they can only rely on different error codes to distinguish different errors.
	// The msg field is used to store additional error information attached when the withErrorCode error is created,
	// The pc field is the program counter when the withErrorCode error was created, which can be used to print the function name + file name + line number when the error was created.
	withErrorCode struct {
		errCode *errorCode
		msg     string
		pc      uintptr
	}
)

func (e *withErrorCode) ErrorCode() *errorCode {
	return e.errCode
}

func (e *withErrorCode) PC() uintptr {
	return e.pc
}

// Error prints errCode.name, errCode.code, errCode.msg, and msg in turn,
// e.g.: ErrUnauthorized, Code=10002, Msg=Unauthorized, something wrong;
// e.g.: ErrUnauthorized, Code=10002, Msg=Unauthorized.
func (e *withErrorCode) Error() string {
	var b strings.Builder

	b.WriteString(e.errCode.name)
	b.WriteString(", Code=")
	b.WriteString(strconv.Itoa(e.errCode.code))
	b.WriteString(", Msg=")
	b.WriteString(e.errCode.msg)

	if e.msg != "" {
		b.WriteString(Config.MessagesSeparator)
		b.WriteString(e.msg)
	}

	return b.String()
}

// Format formats the error message according to the given format specifier.
// It implements the fmt.Formatter interface.
func (e *withErrorCode) Format(s fmt.State, verb rune) {
	formatWithPC(e, s, verb)
}
