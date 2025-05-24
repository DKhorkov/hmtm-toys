package errors

import "fmt"

type MasterNotFoundError struct {
	Message string
	BaseErr error
}

func (e MasterNotFoundError) Error() string {
	template := "master not found"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}

func (e MasterNotFoundError) Unwrap() error {
	return e.BaseErr
}

type MasterAlreadyExistsError struct {
	Message string
	BaseErr error
}

func (e MasterAlreadyExistsError) Error() string {
	template := "master already exists"
	if e.Message != "" {
		template = e.Message
	}

	if e.BaseErr != nil {
		return fmt.Sprintf(template+". Base error: %v", e.BaseErr)
	}

	return template
}

func (e MasterAlreadyExistsError) Unwrap() error {
	return e.BaseErr
}
