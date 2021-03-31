package ppcerrors

import stderrors "errors"

// Wrap 创建一个 withCause 类型的 error，
// err 参数作为底层错误存入 withCause.cause 字段中，
// message 参数用于创建一个 withCaller 类型的错误并存入 withCause.error 字段中作为上层错误，
// 当传入的 err 参数为 nil 时 Wrap 会返回 nil。
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withCause{
		error: &withCaller{
			error: stderrors.New(message),
			pc:    getPCFromCaller(),
		},
		cause: err,
	}
}
