package ppcerrors

import (
	"errors"
	"testing"
)

func TestDefinition(t *testing.T) {
	def := NewDefinition("ErrNotFound", "The requested resource was not found")

	t.Run("Test Name", func(t *testing.T) {
		if def.Name() != "ErrNotFound" {
			t.Errorf("Expected name to be 'ErrNotFound', got '%s'", def.Name())
		}
	})

	t.Run("Test Desc", func(t *testing.T) {
		if def.Desc() != "The requested resource was not found" {
			t.Errorf("Expected description to be 'The requested resource was not found', got '%s'", def.Desc())
		}
	})

	t.Run("Test New", func(t *testing.T) {
		err := def.New("Additional message")
		expectedMsg := "ErrNotFound, The requested resource was not found, Additional message"
		if err.Error() != expectedMsg {
			t.Errorf("Expected error message to be '%s', got '%s'", expectedMsg, err.Error())
		}
	})

	t.Run("Test Wrap", func(t *testing.T) {
		cause := errors.New("Some error")
		err := def.Wrap(cause, "Additional context")
		expectedMsg := "ErrNotFound, The requested resource was not found, Additional context <= Some error"
		if err.Error() != expectedMsg {
			t.Errorf("Expected wrapped error message to be '%s', got '%s'", expectedMsg, err.Error())
		}
	})
}
