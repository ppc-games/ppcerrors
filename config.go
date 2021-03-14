package ppcerrors

// Config 定义了所有可修改配置项。
var Config = struct {
	// 用于在打印日志时分辨是来自于哪个包的配置
	Package string

	// 是否打印创建任何错误时的调用者函数名、所在的文件名、所在行号等信息，默认：false。
	Caller bool `key:"ppc.errors.caller"`
}{
	Package: "ppcerrors",
	Caller:  false,
}
