package ppcerrors

import (
	"fmt"
	"io"
)

// withMessage 实现了 builtin error 接口，并使用 msg 来描述一个错误，
// pc 是 withMessage 错误被创建时的程序计数器，该计数器可用于打印创建错误时执行的函数名+文件名+行号。
type withMessage struct {
	msg string
	pc  uintptr
}

func (e *withMessage) Error() string {
	return e.msg
}

func (e *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// 先打印当前错误
			_, _ = io.WriteString(s, e.Error())

			// 如果记录了 caller，会打印创建错误时调用的函数、文件名、行号
			formatCaller(e.pc, s)

			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}
