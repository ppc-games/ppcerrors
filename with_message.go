package ppcerrors

import (
	"fmt"
)

// withMessage is an error that contains a message and a program counter.
// msg field is used to describe the current error,
// pc field is the program counter when the withMessage error was created, which can be used to print the function name + file name + line number when the error was created.
type withMessage struct {
	msg string
	pc  uintptr
}

func (e *withMessage) PC() uintptr {
	return e.pc
}

// Error returns e.msg.
func (e *withMessage) Error() string {
	return e.msg
}

// Format formats the error message according to the given format specifier.
// It implements the fmt.Formatter interface.
func (e *withMessage) Format(s fmt.State, verb rune) {
	formatWithPC(e, s, verb)
}
