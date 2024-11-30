package errors

type ToyNotFoundError struct {
	Message string
}

func (e ToyNotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "toy not found"
}

type ToyAlreadyExistsError struct {
	Message string
}

func (e ToyAlreadyExistsError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "toy already exists"
}
