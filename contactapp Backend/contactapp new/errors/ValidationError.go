package errors

//Validation represents a validation error
type ValidationError struct {
	Err string `json:"error"`
}

//
func (vError *ValidationError) Error() string {
	return vError.Err
}

func NewValidationError(error string) *ValidationError {
	return &ValidationError{
		Err: error,
	}
}
