package ppcerrors

import (
	"strings"
)

// withDefinition 实现了 builtin error 接口，并使用错误名称 name 和错误描述 desc 描述一个错误。
type withDefinition struct {
	name string
	desc string
}

// NewWithDefinition 创建一个 withDefinition 的实例的指针。
func NewWithDefinition(name string, desc string) *withDefinition {
	return &withDefinition{name: name, desc: desc}
}

// Name 返回当前错误 e 的错误名称。
func (e *withDefinition) Name() string {
	return e.name
}

// Desc 返回当前错误 e 的错误描述。
func (e *withDefinition) Desc() string {
	return e.desc
}

// Error 先打印错误名称，然后打印错误描述，
// 例如：ErrNilUser, 用户信息为空。
func (e *withDefinition) Error() string {
	var b strings.Builder
	b.WriteString(e.name)
	b.WriteString(messagesSeparator)
	b.WriteString(e.desc)
	return b.String()
}

// WithMessage 创建一个 withCause 类型的 error，
// message 参数用于创建一个 withMessage 类型的错误，作为上层错误存入 withCause.error 字段中，
// 当前错误 e 将作为底层错误存入 withCause.cause 中。
func (e *withDefinition) WithMessage(message string) error {
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
func (e *withDefinition) Wrap(cause error, messages ...string) error {
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
