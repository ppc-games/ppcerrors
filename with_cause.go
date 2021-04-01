package ppcerrors

import (
	"fmt"
	"io"
	"strings"
)

// withCause 实现了 builtin error 接口，
// error 字段用于记录上层错误，
// cause 字段用于包装底层错误。
type withCause struct {
	error
	cause error
}

// Error 先打印当前错误的文本描述，然后打印 cause 的文本描述，
// 例如：ErrNilUser, 用户信息为空: cause的Error()。
func (e *withCause) Error() string {
	var b strings.Builder
	b.WriteString(e.error.Error())
	b.WriteString(errorChainSeparator)
	b.WriteString(e.cause.Error())
	return b.String()
}

// Cause 实现了 causer 接口，返回当前错误 e 包装的底层错误 cause。
func (e *withCause) Cause() error {
	return e.cause
}

// Unwrap 实现了 errors 标准库中的 Unwrap 接口，返回当前错误 e 包装的底层错误 cause。
func (e *withCause) Unwrap() error {
	return e.cause
}

// As 用于判断 e.error 是否能赋值给 target。
func (e *withCause) As(target interface{}) bool {
	return As(e.error, target)
}

// Is 用于判断 e.error 是否是一个 err。
func (e *withCause) Is(err error) bool {
	return Is(e.error, err)
}

// Format 当 verb == %+v 时会打印错误链中每一层 cause 的详细错误原因。
func (e *withCause) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// 先打印当前错误
			_, _ = fmt.Fprintf(s, "%+v", e.error)

			// 然后打印被包装的 cause
			_, _ = fmt.Fprintf(s, "\ncause: %+v", e.cause)

			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}
