package ppcerrors

// Define 返回一个错误定义 definition 的实例的指针。
func Define(name string, desc string) *definition {
	return &definition{name: name, desc: desc}
}

// definition 用于定义一个错误，
// name 错误的名称，
// desc 错误的文字描述。
type definition struct {
	name string
	desc string
}

// New 使用当前的错误定义 definition 创建一个 Error 实例。
func (d *definition) New() *Error {
	err := &Error{def: d}

	return err.withCaller()
}

// Wrap 使用当前的错误定义 definition 创建一个 Error 实例，
// 并将原始错误 cause 包装在 Error 中。
func (d *definition) Wrap(cause error) *Error {
	err := &Error{def: d, cause: cause}

	return err.withCaller()
}

// WithErrorCode 使用当前的错误定义 definition 创建一个 Error 实例，
// 并附加一个明确的错误码 errorCode 在 Error 中。
func (d *definition) WithErrorCode(errorCode *errorCode) *Error {
	err := &Error{def: d, errorCode: errorCode}

	return err.withCaller()
}

// WrapWithErrorCode 使用当前的错误定义 definition 创建一个 Error 实例，
// 将原始错误 cause 包装在 Error 中，
// 并附加一个明确的错误码 errorCode 在 Error 中。
func (d *definition) WrapWithErrorCode(cause error, errorCode *errorCode) *Error {
	err := &Error{def: d, cause: cause, errorCode: errorCode}

	return err.withCaller()
}
