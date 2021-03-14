package ppcerrors

type (
	// ErrorCoder 接口定义了 ErrorCode 需要实现的方法。
	ErrorCoder interface {
		Code() int
		Msg() string
	}

	// errorCode 用于定义一个返回给外部系统（例如：客户端）的错误码。
	errorCode struct {
		code int
		msg  string
	}
)

// NewErrorCode 创建并返回一个具体的错误码实例的指针。
func NewErrorCode(code int, msg string) *errorCode {
	return &errorCode{code: code, msg: msg}
}

// Code 返回当前 ErrorCode 的错误码。
func (c *errorCode) Code() int {
	return c.code
}

// Msg 返回当前 ErrorCode 的错误描述。
func (c *errorCode) Msg() string {
	return c.msg
}
