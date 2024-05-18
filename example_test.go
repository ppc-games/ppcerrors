package ppcerrors_test

import (
	"errors"
	"fmt"

	"github.com/ppc-games/ppcerrors"
)

var (
	// Example error definitions
	ErrUpdateOneFailed = ppcerrors.NewDefinition("ErrUpdateOneFailed", "db.UpdateOne failed")
	// Example error codes
	ErrInternalServerError = ppcerrors.NewErrorCode("ErrInternalServerError", 500, "Internal server error")
)

// SaveUser mocks a function that saves a user to the database,
// it will return an error if the database operation fails.
func SaveUser() error {
	err := errors.New("mock mongodb error")
	return ErrUpdateOneFailed.Wrap(err, fmt.Sprintf("SaveUser failed, uid: %d", 123))
}

// Login mocks a function that being invoked when a user requests the login API.
func Login() error {
	if err := SaveUser(); err != nil {
		return ErrInternalServerError.Wrap(err, "Login failed")
	}
	return nil
}

func ExamplePrintWithoutCaller() {
	if err := Login(); err != nil {
		fmt.Println(err)
	}

	// Output: ErrInternalServerError, Code=500, Msg=Internal server error, Login failed <= ErrUpdateOneFailed, db.UpdateOne failed, SaveUser failed, uid: 123 <= mock mongodb error
}

func ExamplePrintWithCaller() {
	ppcerrors.Config.Caller = true

	if err := Login(); err != nil {
		fmt.Printf("%+v", err)
	}

	// Example output:
	// ErrInternalServerError, Code=500, Msg=Internal server error, Login failed
	//     at github.com/ppc-games/ppcerrors_test.Login
	//         /Users/liangrui/Projects/go/ppcerrors/example_test.go:27
	// cause: ErrUpdateOneFailed, db.UpdateOne failed, SaveUser failed, uid: 123
	//     at github.com/ppc-games/ppcerrors_test.SaveUser
	//         /Users/liangrui/Projects/go/ppcerrors/example_test.go:21
	// cause: mock mongodb error
}
