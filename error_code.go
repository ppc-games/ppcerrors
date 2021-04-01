package ppcerrors

import "strings"

type (
	// ErrorCoder 接口定义了 errorCode 需实现的方法。
	ErrorCoder interface {
		Name() string
		Code() int
		Msg() string
	}

	// errorCode 使用 name、code、msg 定义一个错误。
	errorCode struct {
		name string
		code int
		msg  string
	}
)

// NewErrorCode 创建并返回一个错误码 errorCode 实例的指针。
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

// New 基于当前的错误码 c 创建一个 withErrorCode 类型的 error，
// messages 参数用于附加额外的错误信息，使用逗号（,）拼接后存入 msg 字段中，
// 当 Config.Caller == true 时，pc 会记录调用该方法的函数名、文件、行号。
func (c *errorCode) New(messages ...string) error {
	return &withErrorCode{
		errCode: c,
		msg:     strings.Join(messages, messagesSeparator),
		pc:      getPCFromCaller(),
	}
}

// Wrap 基于当前的错误码 c 创建一个 withCause 类型的 error，
// cause 参数将作为底层错误存入 withCause.cause 字段中，
// messages 参数用于附加额外的错误信息，会和当前错误码 c 一起作为上层错误存入 withCause.error 字段中，
// 当传入的 cause 参数为 nil 时 Wrap 会返回 nil。
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
