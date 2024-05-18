package ppcerrors

// Config defines all modifiable configuration items.
var Config = struct {
	// Used to distinguish which package the configuration comes from when printing logs
	Package string
	// Whether to print the caller function name, the file name where it is located, the line number, etc. when any error is created, default: false.
	Caller bool `key:"errors.caller"`
	// Separator connecting two error messages under the same error
	MessagesSeparator string
	// Separator connecting two errors in the error chain
	ErrorChainSeparator string
}{
	Package:             "ppcerrors",
	Caller:              false,
	MessagesSeparator:   ", ",
	ErrorChainSeparator: " <= ",
}
