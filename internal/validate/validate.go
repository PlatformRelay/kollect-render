// Package validate checks inventory documents against schema v0.
package validate

import "fmt"

// Error is a single schema violation located by JSON Pointer.
type Error struct {
	Pointer string
	Message string
}

func (e Error) Error() string {
	if e.Pointer == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Pointer, e.Message)
}

// Result holds zero or more pointered validation errors.
type Result struct {
	Errors []Error
}

func (r Result) Error() string {
	if len(r.Errors) == 0 {
		return ""
	}
	if len(r.Errors) == 1 {
		return r.Errors[0].Error()
	}
	return fmt.Sprintf("%d validation errors (first: %s)", len(r.Errors), r.Errors[0].Error())
}

func (r Result) Err() error {
	if len(r.Errors) == 0 {
		return nil
	}
	return r
}

// Bytes validates a YAML or JSON inventory document.
// Stub: returns a hard failure until schema wiring lands.
func Bytes(_ []byte) error {
	return fmt.Errorf("validate: not implemented")
}

// File reads path and validates it as an inventory document.
func File(_ string) error {
	return fmt.Errorf("validate: not implemented")
}
