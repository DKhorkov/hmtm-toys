package errors

import "fmt"

type ToyNotFoundError struct {
	Message string
	BaseErr error
}

func (e ToyNotFoundError) Error() string {
	template := "toy not found"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}

type ToyAlreadyExistsError struct {
	Message string
	BaseErr error
}

func (e ToyAlreadyExistsError) Error() string {
	template := "toy already exists"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}
