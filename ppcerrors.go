/*
Package ppcerrors enhances the standard errors package and absorbs the design of the github.com/pkg/errors package
to meet the specific error-handling needs of PPC Studio projects,
which makes it easier to identify, append information to, and log errors that are being passed inside and across the systems.

# Print errors wrapped by ppcerrors

By default, the caller information is not printed when an error is created.
For example, the following code wraps an "mock mongodb error" error twice:

	err := errors.New("mock mongodb error")
	err = ErrUpdateOneFailed.Wrap(err, fmt.Sprintf("SaveUser failed, uid: %d", 123))
	err = ErrInternalServerError.Wrap(err, "Login failed")

The default output using fmt.Println(err) will be like this:

	ErrInternalServerError, Code=500, Msg=Internal server error, Login failed <= ErrUpdateOneFailed, db.UpdateOne failed, SaveUser failed, uid: 123 <= mock mongodb error

You will see the error chain in the output orderred from the last wrapped error to the initial error.
However, the output does not contain the function name, package path, file path, and line number of each error's occurrence.

To print the these information, set ppcerrors.Config.Caller to true:

	ppcerrors.Config.Caller = true

Then the output using fmt.Printf("%+v", err) will be like this:

	ErrInternalServerError, Code=500, Msg=Internal server error, Login failed
		at github.com/ppc-games/ppcerrors_test.Login
			/Users/liangrui/Projects/go/ppcerrors/example_test.go:27
	cause: ErrUpdateOneFailed, db.UpdateOne failed, SaveUser failed, uid: 123
		at github.com/ppc-games/ppcerrors_test.SaveUser
			/Users/liangrui/Projects/go/ppcerrors/example_test.go:21
	cause: mock mongodb error

Explanation to the above output:
 1. mock mongodb error is the root cause of the error chain.
 2. ErrUpdateOneFailed is the first error that wraps the root cause inside the SaveUser function.
 3. ErrInternalServerError is the second error that wraps the ErrUpdateOneFailed inside the Login function.

Each error in the chain is printed in three lines:
 1. The error message are printed in the first line.
 2. The function name and package path are printed in the second line.
 3. The file path and line number are printed in the third line.

Now you can easily locate where the error occurred in the code and how the error being passed inside and across the systems.

See more from the [examples](https://pkg.go.dev/github.com/ppc-games/ppcerrors#pkg-examples)

# Use NewDefinition to define errors that normally occur within single services.

For example, defining errors that occur when performing mongo operations:

	var (
		ErrUpdateOneFailed = ppcerrors.NewDefinition("ErrUpdateOneFailed", "db.UpdateOne failed")
		ErrNoDocumentWasUpdated = ppcerrors.NewDefinition("ErrNoDocumentWasUpdated", "No document was updated")
	)

Note: ErrUpdateOneFailed and ErrNoDocumentWasUpdated are type definition structs, not type error.

# Create an actual error type using either definition.New or definition.Wrap.

For example, when an error occurs when saving a user to the database,
you can use the predefined definition to Wrap an existing error with additional information,
or use New to create a new error.

	func SaveUser(user *User) error {
		updateCount, err := db.UpdateOne(user)
		if err != nil {
			return ErrUpdateOneFailed.Wrap(err, fmt.Sprintf("SaveUser failed, uid: %d", user.ID))
		}
		if updateCount == 0 {
			return ErrNoDocumentWasUpdated.New(fmt.Sprintf("SaveUser failed, uid: %d", user.ID))
		}
		return nil
	}

# Identify errors using HasDefinition.

For example, to log MongoDB errors in a middleware:

	func logMongoErrorMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{ResponseWriter: w}
			next.ServeHTTP(rw, r)
			if rw.err != nil {
				switch {
				case ppcerrors.HasDefinition(rw.err, ErrUpdateOneFailed):
					log.Println("db.UpdateOne failed", rw.err)
				case ppcerrors.HasDefinition(rw.err, ErrNoDocumentWasUpdated):
					log.Println("No document was updated", rw.err)
				}
			}
		})
	}

# Use NewErrorCode to define errors that are passed between services.

For example, defining error codes that are returned to the client in HTTP responses:

	var (
		ErrInternalServerError = ppcerrors.NewErrorCode("ErrInternalServerError", 500, "Internal server error")
		ErrUnauthorized = ppcerrors.NewErrorCode("ErrUnauthorized", 401, "Unauthorized")
	)

Note: ErrInternalServerError and ErrUnauthorized are type errorCode structs, not type error.

# Create an actual error type using either errorCode.New or errorCode.Wrap.

For example, when an error occurs when one user requests the login API:

	func Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
		user := ctx.Value("user")
		if user == nil {
			return nil, ErrUnauthorized.New("Login failed, unauthorized user")
		}
		// ... perform login logic ...
		if err:= SaveUser(user); err != nil {
			return nil, ErrInternalServerError.Wrap(err, "Login failed, save user failed")
		}
		return buildLoginResponse(user), nil
	}

Note: different kinds of errors returned by SaveUser are normalized to ErrInternalServerError.

# Identify errors using HasErrorCode.

Refer to the example in the HasDefinition section.

# Always Warp an error to records its initial occurrence even if there is no need to identify the original error later.

For example, wrap an error returns from a third-party library:

	if err := json.Unmarshal(bytesOfState, target); err != nil {
		return ppcerrors.Wrap(err, "json.Unmarshal failed")
	}

Note: There is no need to initialize a definition for each third-party error when no specific error-handling to that error.
For the same reason, there is no need to define an error code if no system cares about that error.
But Wrap is always necessary for the purpose of recording where the error occurred.
*/
package ppcerrors

import (
	stderrors "errors"
)

const (
	// separator connecting two error messages under the same error
	messagesSeparator = ", "
	// separator connecting two errors in the error chain
	errorChainSeparator = " <= "
)

// Wrap creates an error of type withCause.
// The cause parameter is stored in the withCause.cause field as the underlying error,
// and the message parameter is used to create an error of type withMessage, which is stored in the withCause.error field.
// Wrap returns nil when the cause parameter is nil.
func Wrap(cause error, message string) error {
	if cause == nil {
		return nil
	}
	return &withCause{
		error: &withMessage{
			msg: message,
			pc:  getPCFromCaller(),
		},
		cause: cause,
	}
}

// HasErrorCode returns true if err and its error chain contain the specified error code target.
func HasErrorCode(err error, target *errorCode) bool {
	var withErrorCode WithErrorCoder
	if As(err, &withErrorCode) && withErrorCode.ErrorCode() == target {
		return true
	}
	return false
}

// HasDefinition returns true if err and its error chain contain the specified definition target.
func HasDefinition(err error, target *definition) bool {
	var withDefinition WithDefinitioner
	if As(err, &withDefinition) && withDefinition.Definition() == target {
		return true
	}
	return false
}

// Is reports whether any error in err's chain matches the target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool { return stderrors.Is(err, target) }

// As finds the first error in err's chain that matches the target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches the target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool { return stderrors.As(err, target) }

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}
