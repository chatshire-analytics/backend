package errors

type FlipsideError struct {
	Message string
}

func (e *FlipsideError) Error() string {
	return e.Message
}
