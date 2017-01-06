package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// defaultHandler is a http.Handler that provides runtime configuration, request
// logging and panic handling.
type defaultHandler struct {
	Runtime *Runtime
	Handler http.Handler
}

// WrapHandler returns a defaultHandler that wraps the given http.Handler.
func WrapHandler(rt *Runtime, h http.Handler) http.Handler {
	return &defaultHandler{
		Runtime: rt,
		Handler: h,
	}
}

func (c *defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("X-Content-Type-Options", "nosniff")

	defer func() {
		if err := recover(); err != nil {
			status := http.StatusInternalServerError
			if err == MappingNotFoundError {
				status = http.StatusNotFound
			} else if herr, ok := err.(*HTTPError); ok {
				if herr.StatusCode > 0 {
					status = herr.StatusCode
				}
			} else {
				if c.Runtime.Config.ExitOnError {
					panic(err)
				}

				// TODO: log stack traces
			}

			c.Runtime.Logger.Printf("Error: %v", err)

			w.WriteHeader(status)

			if body, err := BodyForStatus(status); err != nil {
				fmt.Fprintf(w, "%d %s\n", status, http.StatusText(status))
			} else {
				fmt.Fprintf(w, body)
			}

		}

		d := time.Since(start)
		c.Runtime.Logger.Printf("%v %v %v", r.Method, r.URL, d)
	}()

	c.Handler.ServeHTTP(w, r)
}

// JSON encodes the given interface{} to JSON and writes the output to the given
// http.ResponseWriter.
func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		panic(err)
	}
}
