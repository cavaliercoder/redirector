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

func loggable(v interface{}) string {
	if v == nil {
		return "-"
	}
	if s, ok := v.(string); ok && s == "" {
		return "-"
	}
	return fmt.Sprintf("%v", v)
}

func (c *defaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ww := NewResponseWriter(w)
	defer func() {
		if err := recover(); err != nil {
			status := StatusCodeForError(err.(error))
			if status >= 500 && status < 600 {
				c.Runtime.Logger.Printf("Error: %v", err)
				// TODO: log stack traces

				if c.Runtime.Config.ExitOnError {
					panic(err)
				}
			}

			ww.WriteHeader(status)
			if body, err := BodyForStatus(status); err != nil {
				c.Runtime.Logger.Printf("Error getting body for status %v: %v", status, err)
				fmt.Fprintf(ww, "%d %s\n", status, http.StatusText(status))
			} else {
				fmt.Fprintf(ww, body)
			}
		}

		d := time.Since(start)
		c.Runtime.AccessLogger.Printf(
			"%s %s %s %s %d %d %s %s %d",
			r.RemoteAddr,
			r.Method,
			r.URL,
			r.Proto,
			ww.Status(),
			ww.Size(),
			loggable(r.Header.Get("Referer")),
			loggable(ww.Header().Get("Location")),
			d.Nanoseconds())
	}()

	ww.Header().Set("Server", PACKAGE_NAME+"/"+PACKAGE_VERSION)
	ww.Header().Set("X-Content-Type-Options", "nosniff")
	c.Handler.ServeHTTP(ww, r)
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
