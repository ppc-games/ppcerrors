# ppcerrors

This library is designed to meet the specific error-handling needs of PPC Studio projects. It offers the following features:

- **Contextual Error Wrapping**: Wrap errors with additional contextual information, such as appending the user ID of the API request initiator when an error occurs.
- **Error Identification**: Use the `HasDefinition` function to compare errors against predefined definitions. Once an error is wrapped with a definition, it can be identified regardless of how many times it is subsequently wrapped.
- **Error Code Handling**: Append error codes to errors for easy identification by external systems. Use the `HasErrorCode` function to detect errors wrapped with specific error codes.
- **Error Normalization**: Normalize different errors (e.g., error A and error B) by wrapping them into the same error (e.g., error C) while preserving the original information of the initial errors. This is useful when errors A and B need to be treated as the same category of error.
- **Detailed Error Reporting**: Record the function name, file name, and line number where the error occurred, and output easy-to-read error reports using built-in print methods. This ensures clear and informative error messages.
- **Efficient Error Stack Printing**: Print the error stack only once, even when the original error is wrapped multiple times.
