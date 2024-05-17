package ppcerrors

import (
	"errors"
	"testing"
)

func TestWrap(t *testing.T) {
	t.Run("Wrap with non-nil cause", func(t *testing.T) {
		cause := errors.New("cause error")
		message := "wrapped error"
		err := Wrap(cause, message)

		if err == nil {
			t.Error("Expected non-nil error, got nil")
		}

		if errWithCause, ok := err.(*withCause); !ok {
			t.Error("Expected wrapped err to be of type withCause")
		} else {
			if errWithCause.Cause() != cause {
				t.Error("Expected wrapped errWithCause.Cause() to return the cause error")
			}

			if errWithMsg, ok := errWithCause.error.(*withMessage); !ok {
				t.Error("Expected wrapped err to be of type withMessage")
			} else {
				if errWithMsg.Error() != message {
					t.Errorf("Expected wrapped errWithMsg.Error() to return '%s', got '%s'", message, errWithMsg.Error())
				}
			}
		}

		if !errors.Is(err, cause) {
			t.Error("Expected wrapped error to be the cause of the error")
		}
	})

	t.Run("Wrap with nil cause", func(t *testing.T) {
		err := Wrap(nil, "wrapped error")
		if err != nil {
			t.Error("Expected nil error, got non-nil")
		}
	})
}

func TestHasErrorCode(t *testing.T) {
	errCode := &errorCode{name: "OK", code: 200, msg: "OK"}

	t.Run("Error with matching error code", func(t *testing.T) {
		err := &withErrorCode{errCode: errCode}
		if !HasErrorCode(err, errCode) {
			t.Error("Expected HasErrorCode to return true")
		}
	})

	t.Run("Error without matching error code", func(t *testing.T) {
		anotherErrCode := &errorCode{name: "NotFound", code: 404, msg: "Not Found"}
		err := &withErrorCode{errCode: anotherErrCode}
		if HasErrorCode(err, errCode) {
			t.Error("Expected HasErrorCode to return false")
		}
	})

	t.Run("Common error", func(t *testing.T) {
		err := errors.New("common error")
		if HasErrorCode(err, errCode) {
			t.Error("Expected HasErrorCode to return false")
		}
	})
}

func TestHasDefinition(t *testing.T) {
	def := &definition{name: "ErrInvalidConfig", desc: "Invalid configuration"}

	t.Run("Error with matching definition", func(t *testing.T) {
		err := &withDefinition{def: def}
		if !HasDefinition(err, def) {
			t.Error("Expected HasDefinition to return true")
		}
	})

	t.Run("Error without matching definition", func(t *testing.T) {
		anotherDef := &definition{name: "ErrInvalidValue", desc: "Invalid value"}
		err := &withDefinition{def: anotherDef}
		if HasDefinition(err, def) {
			t.Error("Expected HasDefinition to return false")
		}
	})

	t.Run("Common error", func(t *testing.T) {
		err := errors.New("common error")
		if HasDefinition(err, def) {
			t.Error("Expected HasDefinition to return false")
		}
	})
}

func TestIs(t *testing.T) {
	target := errors.New("target error")

	t.Run("Is with matching target", func(t *testing.T) {
		err := Wrap(target, "wrapped error")
		if !Is(err, target) {
			t.Error("Expected Is to return true")
		}
	})

	t.Run("Is without matching target", func(t *testing.T) {
		err := errors.New("other error")
		if Is(err, target) {
			t.Error("Expected Is to return false")
		}
	})
}

func TestAs(t *testing.T) {
	var target WithErrorCoder

	err := &withErrorCode{errCode: &errorCode{}}

	t.Run("Error with matching type", func(t *testing.T) {
		if !As(err, &target) {
			t.Error("Expected As to return true")
		} else {
			if target != err {
				t.Error("Expected As to set target to the err")
			}
		}
	})

	t.Run("Wrapped error with matching type", func(t *testing.T) {
		wrappedErr := Wrap(err, "wrapped error")
		if !As(wrappedErr, &target) {
			t.Error("Expected As to return true")
		} else {
			if target != err {
				t.Error("Expected As to set target to the err")
			}
		}
	})

	t.Run("Error without matching type", func(t *testing.T) {
		commonErr := errors.New("common error")
		if As(commonErr, &target) {
			t.Error("Expected As to return false")
		}
	})

	t.Run("Nil error", func(t *testing.T) {
		if As(nil, &target) {
			t.Error("Expected As to return false")
		}
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("Error with non-nil cause", func(t *testing.T) {
		cause := errors.New("cause error")
		err := Wrap(cause, "wrapped error")
		unwrapped := Unwrap(err)
		if unwrapped != cause {
			t.Error("Expected Unwrap to return the cause error")
		}
	})

	t.Run("Error without cause", func(t *testing.T) {
		err := errors.New("error")
		unwrapped := Unwrap(err)
		if unwrapped != nil {
			t.Error("Expected Unwrap to return nil")
		}
	})
}
