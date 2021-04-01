package ppcerrors

import (
	"fmt"
	"strings"
)

type (
	// WithDefinitioner 定义了 Definition() 方法，用于返回 error 所包含的错误定义。
	WithDefinitioner interface {
		Definition() *definition
	}

	// withDefinition 实现了 builtin error 接口，包含一个错误定义 definition 用来区分其它错误，
	// msg 字段用于存储 withDefinition 错误被创建时附加的额外的错误信息，
	// pc 字段是 withDefinition 错误被创建时的程序计数器，该计数器可用于打印创建错误时执行的函数名+文件名+行号。
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

// Error 依次打印 name、desc、msg，
// 例如：ErrNilUser, 用户信息为空，something wrong。
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

func (e *withDefinition) Format(s fmt.State, verb rune) {
	formatWithPC(e, s, verb)
}
