package main

import (
	"fmt"
	"net/http"
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

func StatusCodeForError(err error) int {
	switch err {
	case nil:
		return 0

	case MappingNotFoundError:
		return http.StatusNotFound
	}

	if herr, ok := err.(*HTTPError); ok {
		if herr.StatusCode > 0 {
			return herr.StatusCode
		}
	}

	return http.StatusInternalServerError
}
