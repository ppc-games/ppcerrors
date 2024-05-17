package ppcerrors

import (
	"fmt"
	"io"
	"runtime"

	"github.com/pkg/errors"
)

// getPCFromCaller returns the program counter (PC) when the function is called.
// The PC can be used to print the function name, file name, and line number where the error is created.
// It returns 0 when Config.Caller is set to false.
func getPCFromCaller() uintptr {
	if Config.Caller {
		if pc, _, _, ok := runtime.Caller(2); ok {
			return pc
		}
	}
	return 0
}

// formatWithPC prints the program counter (PC) corresponding to the function, file name, and line number
// when verb == "%+v" and the error contains the program counter (pc).
func formatWithPC(err error, s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// Print the current error first
			_, _ = io.WriteString(s, err.Error())

			// If the error contains the program counter (pc), print the function, file name, and line number
			if pcer, ok := err.(interface{ PC() uintptr }); ok {
				pc := pcer.PC()
				if pc != 0 {
					// Note: This uses the Frame from the github.com/pkg/errors library to format the output.
					// Reference: https://pkg.go.dev/github.com/pkg/errors#Frame.Format
					//
					// Frame formatting verbs:
					// %s    source file
					// %d    source line
					// %n    function name
					// %v    equivalent to %s:%d
					// %+s   function name and path of source file relative to the compile time
					//       GOPATH separated by \n\t (<funcname>\n\t<path>)
					// %+v   equivalent to %+s:%d
					//
					f := errors.Frame(pc)
					// Style 1:
					// at pitaya-multiplayer-games/servers/horserace/handler.(*Handler).Login
					//     /Users/liangrui/Projects/pitaya-horse-race/servers/horserace/handler/login.go:76
					_, _ = fmt.Fprintf(s, "\n    at %+v", f)
					// Style 2:
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
