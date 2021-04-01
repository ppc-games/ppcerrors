package ppcerrors

import (
	stderrors "errors"
	"fmt"
	"io"
	"runtime"

	"github.com/pkg/errors"
)

const (
	// 拼接错误消息的分隔符
	messagesSeparator = ", "
	// 拼接错误链条的分隔符
	errorChainSeparator = " <= "
)

// Wrap 创建一个 withCause 类型的 error，
// cause 参数作为底层错误存入 withCause.cause 字段中，
// message 参数用于创建一个 withMessage 类型的错误，作为上层错误存入 withCause.error 字段中，
// 当传入的 cause 参数为 nil 时 Wrap 会返回 nil。
func Wrap(cause error, message string) error {
	if cause == nil {
		return nil
	}
	return &withCause{
		error: &withMessage{
			msg: message,
			pc:  getPCFromCaller(),
		},
		cause: cause,
	}
}

// HasErrorCode 当 err 以及 err 包装的错误链条中包含指定的错误码 target 时，返回 true。
func HasErrorCode(err error, target *errorCode) bool {
	var withErrorCode WithErrorCoder
	if As(err, &withErrorCode) && withErrorCode.ErrorCode() == target {
		return true
	}
	return false
}

// HasDefinition 当 err 以及 err 包装的错误链条中包含指定的错误定义 target 时，返回 true。
func HasDefinition(err error, target *definition) bool {
	var withDefinition WithDefinitioner
	if As(err, &withDefinition) && withDefinition.Definition() == target {
		return true
	}
	return false
}

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool { return stderrors.Is(err, target) }

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool { return stderrors.As(err, target) }

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// getPCFromCaller 返回函数被调用时的程序计数器，
// 该计数器可用于打印创建错误时执行的函数名+文件名+行号，
// 当 Config.Caller == false 时，会返回 0。
func getPCFromCaller() uintptr {
	if Config.Caller {
		if pc, _, _, ok := runtime.Caller(2); ok {
			return pc
		}
	}
	return 0
}

// formatWithPC 在 verb == "%+v"，且 err 中包含程序计数器 pc 时，会打印程序计数器对应的函数、文件名、行号。
func formatWithPC(err error, s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// 先打印当前错误
			_, _ = io.WriteString(s, err.Error())

			// 如果记录了 caller，会打印创建错误时调用的函数、文件名、行号
			if pcer, ok := err.(interface{ PC() uintptr }); ok {
				pc := pcer.PC()
				if pc != 0 {
					// 注意：这里直接借用 github.com/pkg/errors 库的 Frame 来格式化输出
					// 官方文档参考：https://pkg.go.dev/github.com/pkg/errors#Frame.Format
					//
					// frame 格式化动词参考
					// %s    source file
					// %d    source line
					// %n    function name
					// %v    equivalent to %s:%d
					// %+s   function name and path of source file relative to the compile time
					//       GOPATH separated by \n\t (<funcname>\n\t<path>)
					// %+v   equivalent to %+s:%d
					//
					f := errors.Frame(pc)
					// 样式1：
					// at up-casino-multiplayer-games/servers/horserace/handler.(*Handler).Login
					//     /Users/liangrui/Projects/aig/up-casino/pitaya-horse-race/servers/horserace/handler/login.go:76
					_, _ = fmt.Fprintf(s, "\n    at %+v", f)
					// 样式2：
					// [etcd_service_discovery.go:560/func1()]
					// _, _ = fmt.Fprintf(s, "  at [%s:%d/%n()]\n", f, f, f)
				}
			}

			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, err.Error())
	}
}
