package failure

// Error represents the error model.
type Error struct {
	Code    int    `json:"statusCode"`
	Message string `json:"message"`
}

// NewError creates a new error.
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
