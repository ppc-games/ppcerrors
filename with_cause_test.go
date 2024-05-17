package ppcerrors

import (
	"errors"
	"fmt"
	"testing"
)

func TestWithCause(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		// Create a new error with a cause
		cause := errors.New("root cause")
		err := &withCause{
			error: errors.New("wrapped error"),
			cause: cause,
		}

		expectedError := "wrapped error <= root cause"
		if err.Error() != expectedError {
			t.Errorf("Error() returned %q, expected %q", err.Error(), expectedError)
		}
	})

	t.Run("Cause", func(t *testing.T) {
		cause := errors.New("root cause")
		err := &withCause{cause: cause}

		if err.Cause() != cause {
			t.Errorf("Cause() returned %v, expected %v", err.Cause(), cause)
		}
	})

	t.Run("Unwrap", func(t *testing.T) {
		cause := errors.New("root cause")
		err := &withCause{cause: cause}

		if err.Unwrap() != cause {
			t.Errorf("Unwrap() returned %v, expected %v", err.Unwrap(), cause)
		}
	})

	t.Run("Format", func(t *testing.T) {
		cause := errors.New("root cause")
		err := &withCause{
			error: errors.New("wrapped error"),
			cause: cause,
		}

		// Test Format() method with verb == %+v
		expectedFormat := "wrapped error\ncause: root cause"
		if fmt.Sprintf("%+v", err) != expectedFormat {
			t.Errorf("Format() returned %q, expected %q", fmt.Sprintf("%+v", err), expectedFormat)
		}

		// Test Format() method with verb != %+v
		expectedFormat = "wrapped error <= root cause"
		if fmt.Sprintf("%v", err) != expectedFormat {
			t.Errorf("Format() returned %q, expected %q", fmt.Sprintf("%v", err), expectedFormat)
		}
	})
}
