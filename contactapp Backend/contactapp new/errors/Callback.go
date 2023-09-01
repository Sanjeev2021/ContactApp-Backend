package errors

type CallbackError struct {
	Err string `json:"error"`
}

// Error inerface
func (e CallbackError) Error() string {
	return e.Err
}

// NewCallbackError returns a new CallbackError
func NewCallbackError(error string) *CallbackError {
	return &CallbackError{
		Err: error,
	}
}
