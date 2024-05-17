package ppcerrors

import (
	stderrors "errors"
)

const (
	// separator connecting two error messages under the same error
	messagesSeparator = ", "
	// separator connecting two errors in the error chain
	errorChainSeparator = " <= "
)

// Wrap creates an error of type withCause,
// the cause parameter is stored in the withCause.cause field as the underlying error,
// the message parameter is used to create an error of type withMessage, which is stored in the withCause.error field,
// Wrap returns nil when the cause parameter is nil.
func Wrap(cause error, message string) error {
	if cause == nil {
		return nil
	}
	return &withCause{
		error: &withMessage{
			msg: message,
			pc:  getPCFromCaller(),
		},
		cause: cause,
	}
}

// HasErrorCode returns true if err and its error chain contain the specified error code target.
func HasErrorCode(err error, target *errorCode) bool {
	var withErrorCode WithErrorCoder
	if As(err, &withErrorCode) && withErrorCode.ErrorCode() == target {
		return true
	}
	return false
}

// HasDefinition returns true if err and its error chain contain the specified definition target.
func HasDefinition(err error, target *definition) bool {
	var withDefinition WithDefinitioner
	if As(err, &withDefinition) && withDefinition.Definition() == target {
		return true
	}
	return false
}

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool { return stderrors.Is(err, target) }

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool { return stderrors.As(err, target) }

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}
