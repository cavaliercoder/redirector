package main

import (
	"encoding/json"
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
			status := http.StatusInternalServerError
			if herr, ok := err.(*HTTPError); ok {
				if herr.StatusCode > 0 {
					status = herr.StatusCode
				}
			} else if err == MappingNotFoundError {
				status = http.StatusNotFound
			}

			w.WriteHeader(status)
			fmt.Fprintf(w, "%d %s\n", status, http.StatusText(status))
		}

		d := time.Since(start)
		c.Runtime.Logger.Printf("%v %v %v", r.Method, r.URL, d)
	}()

	c.Handler.ServeHTTP(w, r)
}

func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		panic(err)
	}
}
