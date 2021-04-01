package ppcerrors

import (
	"fmt"
)

// withMessage 实现了 builtin error 接口，
// msg 字段用来描述当前错误，
// pc 字段是 withMessage 错误被创建时的程序计数器，该计数器可用于打印创建错误时执行的函数名+文件名+行号。
type withMessage struct {
	msg string
	pc  uintptr
}

func (e *withMessage) PC() uintptr {
	return e.pc
}

func (e *withMessage) Error() string {
	return e.msg
}

func (e *withMessage) Format(s fmt.State, verb rune) {
	formatWithPC(e, s, verb)
}
