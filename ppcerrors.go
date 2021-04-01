package ppcerrors

const (
	// messagesSeparator 定义拼接错误消息的分隔符
	messagesSeparator = ", "
	// errorChainSeparator 定义拼接错误链条的分隔符
	errorChainSeparator = ": "
)

// Wrap 创建一个 withCause 类型的 error，
// err 参数作为底层错误存入 withCause.cause 字段中，
// message 参数用于创建一个 withMessage 类型的错误，作为上层错误存入 withCause.error 字段中，
// 当传入的 err 参数为 nil 时 Wrap 会返回 nil。
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withCause{
		error: &withMessage{
			msg: message,
			pc:  getPCFromCaller(),
		},
		cause: err,
	}
}
