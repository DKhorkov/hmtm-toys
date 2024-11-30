package errors

type MasterNotFoundError struct {
	Message string
}

func (e MasterNotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "master not found"
}

type MasterAlreadyExistsError struct {
	Message string
}

func (e MasterAlreadyExistsError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "master already exists"
}
