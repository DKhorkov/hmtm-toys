package errors

type CategoryNotFoundError struct {
	Message string
}

func (e CategoryNotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "category not found"
}
