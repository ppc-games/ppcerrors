package ppcerrors

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestWithMessage(t *testing.T) {
	// Set the caller flag to true to make getPCFromCaller() return a valid program counter
	Config.Caller = true

	err := &withMessage{
		msg: "An error occurred",
		pc:  getPCFromCaller(),
	}

	t.Run("Error", func(t *testing.T) {
		expectedError := "An error occurred"
		if err.Error() != expectedError {
			t.Errorf("Error() method returned incorrect error message. Expected: %s, Got: %s", expectedError, err.Error())
		}
	})

	t.Run("Format", func(t *testing.T) {
		// Test %v format
		expectedFormattedError := "An error occurred"
		if fmt.Sprintf("%v", err) != expectedFormattedError {
			t.Errorf("Format() method returned incorrect formatted error message. Expected: %s, Got: %s", expectedFormattedError, fmt.Sprintf("%v", err))
		}

		// Test %+v format
		// Example output:
		// at pitaya-multiplayer-games/servers/horserace/handler.(*Handler).Login
		//     /Users/liangrui/Projects/pitaya-horse-race/servers/horserace/handler/login.go:76
		expectedFormattedError = fmt.Sprintf("%s\n    at %+v", "An error occurred", errors.Frame(err.PC()))
		if fmt.Sprintf("%+v", err) != expectedFormattedError {
			t.Errorf("Format() method returned incorrect formatted error message. Expected: %s, Got: %s", expectedFormattedError, fmt.Sprintf("%+v", err))
		}
	})
}
