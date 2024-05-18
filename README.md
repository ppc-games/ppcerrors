# ppcerrors

[![Go Reference](https://pkg.go.dev/badge/github.com/ppc-games/ppcerrors.svg)](https://pkg.go.dev/github.com/ppc-games/ppcerrors)
[![Go Report Card](https://goreportcard.com/badge/github.com/ppc-games/ppcerrors)](https://goreportcard.com/report/github.com/ppc-games/ppcerrors)

Package ppcerrors enhances the standard errors package and absorbs the design of the github.com/pkg/errors package to meet the specific error-handling needs of PPC Studio projects, which makes it easier to identify, append information to, and log errors that are being passed inside and across the systems.

## Installation

```bash
go get github.com/ppc-games/ppcerrors
```

## Features

- **Contextual Error Wrapping**: Wrap errors with additional contextual information, such as appending the user ID of the API request initiator when an error occurs.
- **Error Identification**: Use the `HasDefinition` function to compare errors against predefined definitions. Once an error is wrapped with a definition, it can be identified regardless of how many times it is subsequently wrapped.
- **Error Code Handling**: Append error codes to errors for easy identification by external systems. Use the `HasErrorCode` function to detect errors wrapped with specific error codes.
- **Error Normalization**: Normalize different errors (e.g., error A and error B) by wrapping them into the same error (e.g., error C) while preserving the original information of the initial errors. This is useful when errors A and B need to be treated as the same category of error.
- **Detailed Error Reporting**: Record the function name, file name, and line number where the error occurred, and output easy-to-read error reports using built-in print methods. This ensures clear and informative error messages.
- **Efficient Error Stack Printing**: Print the error stack only once, even when the original error is wrapped multiple times.

## Print errors wrapped by ppcerrors

By default, the caller information is not printed when an error is created.
For example, the following code wraps an "mock mongodb error" error twice:

```go
err := errors.New("mock mongodb error")
err = ErrUpdateOneFailed.Wrap(err, fmt.Sprintf("SaveUser failed, uid: %d", 123))
err = ErrInternalServerError.Wrap(err, "Login failed")
```

The default output using fmt.Println(err) will be like this:

`ErrInternalServerError, Code=500, Msg=Internal server error, Login failed <= ErrUpdateOneFailed, db.UpdateOne failed, SaveUser failed, uid: 123 <= mock mongodb error`

You will see the error chain in the output orderred from the last wrapped error to the initial error.
However, the output does not contain the function name, package path, file path, and line number where each error occurred.

To print the these information, set ppcerrors.Config.Caller to true:

```go
ppcerrors.Config.Caller = true
```

Then the output using fmt.Printf("%+v", err) will be like this:

```
ErrInternalServerError, Code=500, Msg=Internal server error, Login failed
    at github.com/ppc-games/ppcerrors_test.Login
        /Users/liangrui/Projects/go/ppcerrors/example_test.go:27
cause: ErrUpdateOneFailed, db.UpdateOne failed, SaveUser failed, uid: 123
    at github.com/ppc-games/ppcerrors_test.SaveUser
        /Users/liangrui/Projects/go/ppcerrors/example_test.go:21
cause: mock mongodb error
```

Explanation to the above output:

1.  mock mongodb error is the root cause of the error chain.
2.  ErrUpdateOneFailed is the first error that wraps the root cause inside the SaveUser function.
3.  ErrInternalServerError is the second error that wraps the ErrUpdateOneFailed inside the Login function.

Each error in the chain is printed in three lines:

1.  The error message are printed in the first line.
2.  The function name and package path are printed in the second line.
3.  The file path and line number are printed in the third line.

Now you can easily locate where the error occurred in the code and how the error being passed inside and across the systems.

See the [Documentation](https://pkg.go.dev/github.com/ppc-games/ppcerrors#section-documentation) for more details.
