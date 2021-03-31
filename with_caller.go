package ppcerrors

import (
	"fmt"
	"io"
	"runtime"

	"github.com/pkg/errors"
)

// withCaller 实现了 builtin error 接口，error 是其包装的原始错误，
// pc 是 withCaller 错误被创建时的程序计数器，该计数器可用于打印创建错误时执行的函数名+文件名+行号。
type withCaller struct {
	error
	pc uintptr
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

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *withCaller) Unwrap() error {
	return e.error
}

// Format 当 verb == %+v 时会打印创建错误时调用的函数、文件名、行号。
func (e *withCaller) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// 先打印当前错误
			_, _ = io.WriteString(s, e.Error())

			// 如果记录了 caller，会打印创建错误时调用的函数、文件名、行号
			if e.pc != 0 {
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
				f := errors.Frame(e.pc)
				// 样式1：
				// at up-casino-multiplayer-games/servers/horserace/handler.(*Handler).Login
				//     /Users/liangrui/Projects/aig/up-casino/pitaya-horse-race/servers/horserace/handler/login.go:76
				_, _ = fmt.Fprintf(s, "\n    at %+v", f)
				// 样式2：
				// [etcd_service_discovery.go:560/func1()]
				// _, _ = fmt.Fprintf(s, "  at [%s:%d/%n()]\n", f, f, f)
			}
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}
