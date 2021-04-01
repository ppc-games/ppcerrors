package ppcerrors

import "strings"

// definition 使用 name 和 desc 定义一个错误。
type definition struct {
	name string
	desc string
}

// NewDefinition 创建并返回一个错误定义 definition 实例的指针。
func NewDefinition(name string, desc string) *definition {
	return &definition{
		name: name,
		desc: desc,
	}
}

func (d *definition) Name() string {
	return d.name
}

func (d *definition) Desc() string {
	return d.desc
}

// New 基于当前的错误定义 d 创建一个 withDefinition 类型的 error，
// messages 参数用于附加额外的错误信息，使用逗号（,）拼接后存入 msg 字段中，
// 当 Config.Caller == true 时，pc 会记录调用该方法的函数名、文件、行号。
func (d *definition) New(messages ...string) error {
	return &withDefinition{
		def: d,
		msg: strings.Join(messages, messagesSeparator),
		pc:  getPCFromCaller(),
	}
}

// Wrap 基于当前的错误定义 d 创建一个 withCause 类型的 error，
// cause 参数将作为底层错误存入 withCause.cause 字段中，
// messages 参数用于附加额外的错误信息，会和当前错误定义 d 一起作为上层错误存入 withCause.error 字段中，
// 当传入的 cause 参数为 nil 时 Wrap 会返回 nil。
func (d *definition) Wrap(cause error, messages ...string) error {
	if cause == nil {
		return nil
	}

	return &withCause{
		error: &withDefinition{
			def: d,
			msg: strings.Join(messages, messagesSeparator),
			pc:  getPCFromCaller(),
		},
		cause: cause,
	}
}
