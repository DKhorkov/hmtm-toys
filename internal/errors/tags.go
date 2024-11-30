package errors

type TagNotFoundError struct {
	Message string
}

func (e TagNotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "tag not found"
}
