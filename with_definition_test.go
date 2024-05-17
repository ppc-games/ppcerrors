package ppcerrors

import (
	"testing"
)

func TestWithDefinition(t *testing.T) {
	def := &definition{
		name: "ErrInvalidConfig",
		desc: "Invalid configuration",
	}
	err := &withDefinition{
		def: def,
		msg: "Test error message",
	}

	expected := def
	actual := err.Definition()

	if actual != expected {
		t.Errorf("Expected definition '%v', but got '%v'", expected, actual)
	}

	expectedError := "ErrInvalidConfig, Invalid configuration, Test error message"
	actualError := err.Error()

	if actualError != expectedError {
		t.Errorf("Expected error message '%s', but got '%s'", expectedError, actualError)
	}
}
