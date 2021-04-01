package ppcerrors

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	// WithErrorCoder 定义了 ErrorCode() 方法，用于返回 error 所包含的错误码。
	WithErrorCoder interface {
		ErrorCode() *errorCode
	}

	// withErrorCode 实现了 builtin error 接口，包含一个错误码 errorCode 用来区分其它错误，
	// 使用错误码的方式一般用于返回给当前应用程序系统边界之外的系统（例如：客户端），
	// 因为这些外部系统无法直接拿到 error，只能依靠识别不同的错误码来区分不同的错误。
	// msg 字段用于存储 withErrorCode 错误被创建时附加的额外的错误信息，
	// pc 字段是 withErrorCode 错误被创建时的程序计数器，该计数器可用于打印创建错误时执行的函数名+文件名+行号。
	withErrorCode struct {
		errCode *errorCode
		msg     string
		pc      uintptr
	}
)

func (e *withErrorCode) ErrorCode() *errorCode {
	return e.errCode
}

func (e *withErrorCode) PC() uintptr {
	return e.pc
}

// Error 依次打印 errCode.name、errCode.code、errCode.msg、msg，
// 例如：ErrUnauthorized, Code=10002, Msg=未授权, something wrong；
// 例如：ErrUnauthorized, Code=10002, Msg=未授权。
func (e *withErrorCode) Error() string {
	var b strings.Builder

	b.WriteString(e.errCode.name)
	b.WriteString(", Code=")
	b.WriteString(strconv.Itoa(e.errCode.code))
	b.WriteString(", Msg=")
	b.WriteString(e.errCode.msg)

	if e.msg != "" {
		b.WriteString(messagesSeparator)
		b.WriteString(e.msg)
	}

	return b.String()
}

func (e *withErrorCode) Format(s fmt.State, verb rune) {
	formatWithPC(e, s, verb)
}
