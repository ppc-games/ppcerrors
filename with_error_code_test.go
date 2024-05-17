package ppcerrors

import (
	"testing"
)

func TestWithErrorCode(t *testing.T) {
	errCode := &errorCode{
		name: "ErrUnauthorized",
		code: 401,
		msg:  "Unauthorized",
	}

	t.Run("Error message should contain error code, message, and additional message", func(t *testing.T) {
		err := &withErrorCode{errCode: errCode, msg: "something wrong"}
		expected := "ErrUnauthorized, Code=401, Msg=Unauthorized, something wrong"
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected error message '%s', but got '%s'", expected, actual)
		}
	})

	t.Run("Error message should contain error code and message", func(t *testing.T) {
		err := &withErrorCode{errCode: errCode}
		expected := "ErrUnauthorized, Code=401, Msg=Unauthorized"
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected error message '%s', but got '%s'", expected, actual)
		}
	})
}
