package ppcerrors

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Error 是一个 error 类型的错误。
// cause 存放包装的原始错误，未包装原始错误时为 nil；
// def 存放错误的名称和文字描述；
// *errorCode 是附加的错误码，未附加错误码时未 nil，
// pc 存放错误被创建时调用函数的程序计数器
type Error struct {
	cause error
	def   *definition // 注意：这里不要嵌入，因为不想暴露 definition 的 New() 和 Wrap() 方法
	*errorCode
	pc uintptr
}

// Definition 返回当前 Error 对应的错误定义，包括错误的名称和文字描述。
func (e *Error) Definition() *definition {
	return e.def
}

// Cause 实现了 github.com/pkg/errors 中的 causer 接口，返回当前 Error 包装的原始错误。
func (e *Error) Cause() error {
	return e.cause
}

// Unwrap 实现了 errors 标准库中的 Unwrap 接口，返回当前 Error 包装的原始错误。
func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) withCaller() *Error {
	// 只有在配置开启需要记录调用者名称时，才进行记录
	if Config.Caller {
		if pc, _, _, ok := runtime.Caller(2); ok {
			e.pc = pc
		}
	}
	return e
}

// Format 模仿 github.com/pkg/errors 的实现方式，识别不同的 fmt 格式化动词将错误内容格式化为字符串。
// 注意：针对 %+v 会打印多层包装中每一层的详细错误原因，
// 如果发现被包装的错误记录了错误堆栈，会打印被包装的错误被创建时的函数名+文件名+行号，但是不会打印整个调用堆栈，
// 如果被包装的错误没有记录错误堆栈，那么会递归地打印 cause 的 Error()。
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// 打印当前层的包装原因
			_, _ = io.WriteString(s, e.ErrorWithoutCause())
			_, _ = io.WriteString(s, "\n")

			// 如果发现有记录调用者，那么打印执行错误包装时的调用函数，文件名，函数
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
				_, _ = fmt.Fprintf(s, "    at %+v\n", f)
				// 样式2：
				// [etcd_service_discovery.go:560/func1()]
				//_, _ = fmt.Fprintf(s, "  at [%s:%d/%n()]\n", f, f, f)
			}

			// 打印被包装的错误的详细错误原因
			if e.cause != nil {
				_, _ = fmt.Fprintf(s, "cause: %+v\n", e.cause)
			}
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}

// Error 返回错误的字符串文本描述，
// 例如：ErrNilUser, 用户信息为空；
// 例如：ErrNilUser, 用户信息为空, Code=10002, Msg=未授权；
// 例如：ErrNilUser, 用户信息为空, Code=10002, Msg=未授权: cause的Error()
func (e *Error) Error() string {
	var b strings.Builder
	e.writeDefinition(&b)
	e.writeErrorCode(&b)
	e.writeCause(&b)
	return b.String()
}

// Error 仅返回当前错误的字符串文本描述，不会增加所包装的 cause 的错误的文本描述。
func (e *Error) ErrorWithoutCause() string {
	var b strings.Builder
	e.writeDefinition(&b)
	e.writeErrorCode(&b)
	return b.String()
}

func (e *Error) writeDefinition(b *strings.Builder) {
	b.WriteString(e.def.name)
	b.WriteString(", ")
	b.WriteString(e.def.desc)
}

func (e *Error) writeErrorCode(b *strings.Builder) {
	if e.errorCode != nil {
		b.WriteString(", Code=")
		b.WriteString(strconv.Itoa(e.errorCode.code))
		b.WriteString(", Msg=")
		b.WriteString(e.errorCode.msg)
	}
}

func (e *Error) writeCause(b *strings.Builder) {
	if e.cause != nil {
		b.WriteString(": ")
		b.WriteString(e.cause.Error())
	}
}
