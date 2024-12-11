package errors

import "fmt"

type TagNotFoundError struct {
	Message string
	BaseErr error
}

func (e TagNotFoundError) Error() string {
	template := "tag not found"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}
