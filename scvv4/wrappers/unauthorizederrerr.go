package wrappers

// UnauthenticatedErr is an error of type unauthorizedError with the underlying error message
var UnauthenticatedErr error = unauthenticatedError{msg: "aunauthenticated"}

// unauthenticatedError is an implementation of error interface
type unauthenticatedError struct {
	msg string
}

// NewUnauthenticatedErr wraps the given error in an unauthenticatedError
func NewUnauthenticatedErr(err error) error {
	if err == nil {
		return nil
	}

	return unauthenticatedError{
		msg: err.Error(),
	}
}

// Error returns the error message
func (e unauthenticatedError) Error() string {
	return e.msg
}

// Is returns true if the target error is an unauthenticatedError
func (e unauthenticatedError) Is(tgt error) bool {
	_, ok := tgt.(unauthenticatedError)
	return ok
}
