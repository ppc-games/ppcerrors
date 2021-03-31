package ppcerrors

import (
	stderrors "errors"
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

// Error 返回包含错误定义的文本描述，例如：ErrNilUser, 用户信息为空。
func (e *withDefinition) Error() string {
	var b strings.Builder
	b.WriteString(e.name)
	b.WriteString(", ")
	b.WriteString(e.desc)
	return b.String()
}

// WithMessage 返回一个新的 error，其值为 withCause 实例的指针，
// 当前错误 e 将作为底层错误存入 withCause.cause 中，
// message 参数用于创建一个 withCaller 类型的错误并存入 withCause.error 字段中作为上层错误。
func (e *withDefinition) WithMessage(message string) error {
	return &withCause{
		error: &withCaller{
			error: stderrors.New(message),
			pc:    getPCFromCaller(),
		},
		cause: e,
	}
}

// Wrap 创建一个 error，使用当前的错误 e 作为底层错误，并将另一个错误 cause 包装在内部，
// message 用于附加一段错误描述文字。
func (e *withDefinition) Wrap(cause error, messages ...string) error {
	if len(messages) > 0 {
		return &withCause{
			error: &withCause{
				error: &withCaller{
					error: stderrors.New(strings.Join(messages, ", ")),
					pc:    getPCFromCaller(),
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

// Name 返回当前错误 e 的错误名称。
func (e *withDefinition) Name() string {
	return e.name
}

// Desc 返回当前错误 e 的错误描述。
func (e *withDefinition) Desc() string {
	return e.desc
}
