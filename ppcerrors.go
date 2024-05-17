package ppcerrors

import (
	stderrors "errors"
)

const (
	// 拼接错误消息的分隔符
	messagesSeparator = ", "
	// 拼接错误链条的分隔符
	errorChainSeparator = " <= "
)

// Wrap 创建一个 withCause 类型的 error，
// cause 参数作为底层错误存入 withCause.cause 字段中，
// message 参数用于创建一个 withMessage 类型的错误，作为上层错误存入 withCause.error 字段中，
// 当传入的 cause 参数为 nil 时 Wrap 会返回 nil。
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

// HasErrorCode 当 err 以及 err 包装的错误链条中包含指定的错误码 target 时，返回 true。
func HasErrorCode(err error, target *errorCode) bool {
	var withErrorCode WithErrorCoder
	if As(err, &withErrorCode) && withErrorCode.ErrorCode() == target {
		return true
	}
	return false
}

// HasDefinition 当 err 以及 err 包装的错误链条中包含指定的错误定义 target 时，返回 true。
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
