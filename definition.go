package ppcerrors

import "strings"

// definition defines an error with a name and description.
// name is the name of the definition, eg: "ErrNotFound".
// desc is the description of the definition, eg: "The requested resource was not found".
type definition struct {
	name string
	desc string
}

// NewDefinition creates and returns a pointer to an error definition instance.
func NewDefinition(name string, desc string) *definition {
	return &definition{
		name: name,
		desc: desc,
	}
}

func (d *definition) Name() string {
	return d.name
}

func (d *definition) Desc() string {
	return d.desc
}

// New creates a withDefinition error based on the current error definition d,
// the messages parameter is used to attach additional error information,
// which is concatenated with the value of messagesSeparator and stored in the msg field,
// when Config.Caller == true, pc records the function name, file, and line number of the method that called this method.
func (d *definition) New(messages ...string) error {
	return &withDefinition{
		def: d,
		msg: strings.Join(messages, messagesSeparator),
		pc:  getPCFromCaller(),
	}
}

// Wrap wraps the given error with additional context and returns a new error.
// If the cause error is nil, it returns nil.
// The additional context is specified by the messages parameter, which is joined
// using the messagesSeparator.
// The returned error contains the original error, the definition, the joined messages,
// and the program counter of the caller.
func (d *definition) Wrap(cause error, messages ...string) error {
	if cause == nil {
		return nil
	}

	return &withCause{
		error: &withDefinition{
			def: d,
			msg: strings.Join(messages, messagesSeparator),
			pc:  getPCFromCaller(),
		},
		cause: cause,
	}
}
