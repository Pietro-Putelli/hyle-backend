package failure

//-------------------------------------
// Validation
//-------------------------------------

// ValidationErr represents a validation error.
type ValidationErr struct {
	Err error
}

// NewValidationErr creates a new validation error.
func NewValidationErr(err error) *ValidationErr {
	return &ValidationErr{Err: err}
}

// Error returns the error message.
func (e *ValidationErr) Error() string {
	return e.Err.Error()
}
