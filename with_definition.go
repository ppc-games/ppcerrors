package ppcerrors

import (
	"fmt"
	"strings"
)

type (
	// WithDefinitioner defines the Definition() method to return the error definition contained in the error.
	WithDefinitioner interface {
		Definition() *definition
	}

	// withDefinition is an error that contains a definition to distinguish it from other errors.
	// The msg field is used to store additional error information attached when the withDefinition error is created,
	// The pc field is the program counter when the withDefinition error was created, which can be used to print the function name + file name + line number when the error was created.
	withDefinition struct {
		def *definition
		msg string
		pc  uintptr
	}
)

func (e *withDefinition) Definition() *definition {
	return e.def
}

func (e *withDefinition) PC() uintptr {
	return e.pc
}

// Error prints name, desc, and msg in turn,
// e.g.: ErrNilUser, User information is empty, something wrong.
func (e *withDefinition) Error() string {
	var b strings.Builder

	b.WriteString(e.def.name)
	b.WriteString(messagesSeparator)
	b.WriteString(e.def.desc)

	if e.msg != "" {
		b.WriteString(messagesSeparator)
		b.WriteString(e.msg)
	}

	return b.String()
}

// Format formats the error message according to the given format specifier.
// It implements the fmt.Formatter interface.
func (e *withDefinition) Format(s fmt.State, verb rune) {
	formatWithPC(e, s, verb)
}
