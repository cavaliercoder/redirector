package main

import (
	"fmt"
	"net/http"
)

// getMappingOrDefault returns the requested mapping or the mapping for the
// default key or MappingNotFoundError if neither are found.
func getMappingOrDefault(rt *Runtime, key string) (*Mapping, error) {
	keys := make([]string, 0, 2)

	if key != "" {
		keys = append(keys, key)
	}

	if rt.Config.DefaultKey != "" {
		keys = append(keys, rt.Config.DefaultKey)
	}

	for _, k := range keys {
		m, err := rt.Database.GetMapping(k)
		if err == nil {
			if err := m.Validate(); err != nil {
				return nil, err
			}
			return m, nil
		}

		if err != MappingNotFoundError {
			panic(err)
		}
	}

	return nil, NewHTTPError(http.StatusNotFound, MappingNotFoundError)
}

func RedirectHandler(rt *Runtime) http.Handler {
	return WrapHandler(rt, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, err := rt.Config.KeyBuilder.Parse(r)
		if err != nil {
			if status := StatusCodeForError(err); status >= 500 {
				panic(err)
			}
		}

		m, err := getMappingOrDefault(rt, key)
		if err != nil {
			panic(err)
		}

		status := http.StatusTemporaryRedirect
		if m.Permanent {
			status = http.StatusPermanentRedirect
		}

		dest := m.Destination
		if m.IsTemplate {
			if d, err := m.ComputeDestination(key, r); err != nil {
				panic(err)
			} else {
				dest = d
			}
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Location", dest)
		w.WriteHeader(status)

		body, err := BodyForStatus(status)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, body)
	}))
}

func serve(rt *Runtime) error {
	s := &http.Server{
		Addr:    rt.Config.ListenAddr,
		Handler: RedirectHandler(rt),
	}

	rt.Logger.Printf("Listening for redirect requests on %v", rt.Config.ListenAddr)
	return s.ListenAndServe()
}
