package ppcerrors

import (
	"errors"
	"testing"
)

func TestErrorCode(t *testing.T) {
	t.Run("TestNewErrorCode", func(t *testing.T) {
		name := "ErrInternalServerError"
		code := 500
		msg := "Internal server error"
		errCode := NewErrorCode(name, code, msg)

		if errCode.Name() != name {
			t.Errorf("Expected name to be %s, but got %s", name, errCode.Name())
		}

		if errCode.Code() != code {
			t.Errorf("Expected code to be %d, but got %d", code, errCode.Code())
		}

		if errCode.Msg() != msg {
			t.Errorf("Expected message to be %s, but got %s", msg, errCode.Msg())
		}
	})

	t.Run("TestNew", func(t *testing.T) {
		errCode := NewErrorCode("ErrInternalServerError", 500, "Internal server error")
		err := errCode.New("Additional message")
		expectedMsg := "ErrInternalServerError, Code=500, Msg=Internal server error, Additional message"
		if err.Error() != expectedMsg {
			t.Errorf("Expected error message to be %s, but got %s", expectedMsg, err.Error())
		}
	})

	t.Run("TestWrap", func(t *testing.T) {
		errCode := NewErrorCode("ErrInternalServerError", 500, "Internal server error")
		cause := errors.New("This is the cause error")
		err := errCode.Wrap(cause, "Additional message")
		expectedMsg := "ErrInternalServerError, Code=500, Msg=Internal server error, Additional message <= This is the cause error"
		if err.Error() != expectedMsg {
			t.Errorf("Expected error message to be %s, but got %s", expectedMsg, err.Error())
		}
	})
}
