package ppcerrors

import (
	"fmt"
	"io"
	"strings"
)

// withCause 是一个使用 error 字段记录底层错误，并将另一个错误 cause 包装在内部，并实现了 builtin error 接口。
type withCause struct {
	error
	cause error
}

// Error 返回包含所包装错误 cause 的字符串描述，
// 先打印当前错误的文本描述，然后打印 cause 的文本描述，
// 例如：ErrNilUser, 用户信息为空: cause的Error()。
func (e *withCause) Error() string {
	var b strings.Builder
	b.WriteString(e.error.Error())
	b.WriteString(": ") // 多个错误链条使用冒号（:）拼接
	b.WriteString(e.cause.Error())
	return b.String()
}

// Cause 实现了 causer 接口，返回当前错误 e 包装的 cause。
func (e *withCause) Cause() error {
	return e.cause
}

// Unwrap 实现了 errors 标准库中的 Unwrap 接口，返回当前错误 e 包装的 cause。
func (e *withCause) Unwrap() error {
	return e.cause
}

// As 针对 e.error 执行 errors.As 方法，
// 当执行 errors.As(e, target) 时，
// 会先判断 e.error 是否符合 target，然后再判断 e.cause 是否符合 target。
func (e *withCause) As(target interface{}) bool {
	return As(e.error, target)
}

// Is 针对 e.error 执行 errors.Is 方法，
// 当执行 errors.Is(e, err) 时，
// 会先判断 e.error 是否是 err，然后再判断 e.cause 是否是 err。
func (e *withCause) Is(err error) bool {
	return Is(e.error, err)
}

// Format 当 verb == %+v 时会打印多层包装中每一层 cause 的详细错误原因。
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
