package errox

import (
	"fmt"
)

type roxError struct {
	message string
	base    error
}

// makeSentinel returns a new sentinel error. Semantically this is very close to
// `errors.New(message)` from the standard library.
func makeSentinel(message string) error {
	return &roxError{message, nil}
}

// Error returns error message. Implements error interface.
func (e *roxError) Error() string {
	return e.message
}

// Unwrap returns the base of the error.
func (e *roxError) Unwrap() error {
	return e.base
}

// New creates an error based on another error, e.g., an existing sentinel
// error, but with the personalized error message. Essentially, it allows to
// preserve the error base error in the chain but hide its message.
//
// Example:
//     myPackageSentinel := errox.New(errox.NotFound, "gone")
//     myPackageSentinel.Error() == "gone" // true
func New(base error, message string) error {
	return &roxError{message, base}
}

// Newf creates an error based on another error, e.g., an existing sentinel
// error, but with the personalized formatted error message. Essentially, it
// allows to preserve the error base error in the chain but hide its message.
//
// Example:
//     myPackageError := errox.Newf(errox.NotFound, "gone: %s", reason)
//     myPackageError.Error() == "gone: <reason>" // true
func Newf(base error, format string, args ...interface{}) error {
	return New(base, fmt.Sprintf(format, args...))
}