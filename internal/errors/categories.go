package errors

import "fmt"

type CategoryNotFoundError struct {
	Message string
	BaseErr error
}

func (e CategoryNotFoundError) Error() string {
	template := "category not found"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}

func (e CategoryNotFoundError) Unwrap() error {
	return e.BaseErr
}
