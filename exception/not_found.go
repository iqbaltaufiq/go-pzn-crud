package exception

// Go doesn't have built-in http error handling
// make your own error handler
// by implementing Error interface from Go

// Handle error when no item is returned by database
type NotFoundError struct {
	Error string
}

func NewNotFoundError(err string) NotFoundError {
	return NotFoundError{Error: err}
}
