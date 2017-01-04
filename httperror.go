package main

import (
	"fmt"
)

// HTTPError wraps a regular error and provides HTTP Status Code information.
type HTTPError struct {
	Err        error
	StatusCode int
}

func (e *HTTPError) Error() string {
	return e.Err.Error()
}

func NewHTTPError(statusCode int, err error) error {
	return &HTTPError{
		Err:        err,
		StatusCode: statusCode,
	}
}

func NewHTTPErrorf(statusCode int, format string, a ...interface{}) error {
	return NewHTTPError(statusCode, fmt.Errorf(format, a...))
}
