package ppcerrors

import (
	"strconv"
	"strings"
)

type (
	// ErrorCoder 接口定义了 withErrorCode 实现的方法，可用于识别一个 error 是否包含错误码。
	ErrorCoder interface {
		Name() string
		Code() int
		Msg() string
	}

	// withErrorCode 实现了 builtin error 接口，并使用错误名称 name、错误码 code 和错误描述 msg 来描述一个错误，
	// 使用错误码的方式一般用于返回给当前应用程序系统边界之外的系统（例如：客户端），
	// 因为这些外部系统无法直接拿到 error，只能依靠识别不同的错误码来区分不同的错误。
	withErrorCode struct {
		name string
		code int
		msg  string
	}
)

// NewWithErrorCode 创建并返回一个 withErrorCode 实例的指针。
func NewWithErrorCode(name string, code int, msg string) *withErrorCode {
	return &withErrorCode{name: name, code: code, msg: msg}
}

// Error 先打印错误名称，然后打印错误码，最后打印错误描述，
// 例如：ErrUnauthorized, Code=10002, Msg=未授权。
func (e *withErrorCode) Error() string {
	var b strings.Builder
	b.WriteString(e.name)
	b.WriteString(", Code=")
	b.WriteString(strconv.Itoa(e.code))
	b.WriteString(", Msg=")
	b.WriteString(e.msg)
	return b.String()
}

// Name 返回当前错误 e 的错误名称。
func (e *withErrorCode) Name() string {
	return e.name
}

// Code 返回当前错误 e 的错误码。
func (e *withErrorCode) Code() int {
	return e.code
}

// Msg 返回当前错误 e 的错误描述。
func (e *withErrorCode) Msg() string {
	return e.msg
}

// WithMessage 创建一个 withCause 类型的 error，
// message 参数用于创建一个 withMessage 类型的错误，作为上层错误存入 withCause.error 字段中，
// 当前错误 e 将作为底层错误存入 withCause.cause 中。
func (e *withErrorCode) WithMessage(message string) error {
	return &withCause{
		error: &withMessage{
			msg: message,
			pc:  getPCFromCaller(),
		},
		cause: e,
	}
}

// Wrap 创建一个 withCause 类型的 error，
// cause 参数作为底层错误存入 withCause.cause 字段中，
// messages 参数用于创建一个 withMessage 类型的错误，和当前错误 e 一起作为上层错误存入 withCause.error 字段中，
// 当传入的 cause 参数为 nil 时 Wrap 会返回 nil。
func (e *withErrorCode) Wrap(cause error, messages ...string) error {
	if cause == nil {
		return nil
	}

	if len(messages) > 0 {
		return &withCause{
			error: &withCause{
				error: &withMessage{
					msg: strings.Join(messages, messagesSeparator),
					pc:  getPCFromCaller(),
				},
				cause: e,
			},
			cause: cause,
		}
	}

	return &withCause{
		error: &withCaller{
			error: e,
			pc:    getPCFromCaller(),
		},
		cause: cause,
	}
}
