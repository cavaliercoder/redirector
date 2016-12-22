package main

import (
	"fmt"
	"net/http"
	"time"
)

type defaultHandler struct {
	Runtime *Runtime
	Handler http.Handler
}

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
			if herr, ok := err.(*HTTPError); ok {
				status := herr.StatusCode
				if status == 0 {
					status = http.StatusInternalServerError
				}

				w.WriteHeader(status)
				fmt.Fprintf(w, "%d %s\n", status, http.StatusText(status))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%d %s\n", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
		}

		d := time.Since(start)
		c.Runtime.Logger.Printf("%v %v %v", r.Method, r.URL, d)
	}()

	c.Handler.ServeHTTP(w, r)
}
