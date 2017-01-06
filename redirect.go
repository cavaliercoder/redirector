package main

import (
	"fmt"
	"net/http"
)

const (
	// TODO: Update title as per actual response code
	REDIRECT_BODY = `<html>
<head><title>301 Moved Permanently</title></head>
<body bgcolor="white">
<center><h1>301 Moved Permanently</h1></center>
<hr><center>` + PACKAGE_NAME + `/` + PACKAGE_VERSION + `</center>
</body>
</html>`
)

// getMappingOrDefault returns the requested mapping or the mapping for the
// default key or MappingNotFoundError if neither are found.
func getMappingOrDefault(rt *Runtime, key string) (*Mapping, error) {
	keys := []string{key}
	if rt.Config.DefaultKey != "" {
		keys = append(keys, rt.Config.DefaultKey)
	}

	for _, k := range keys {
		m, err := rt.Database.GetMapping(k)
		if err == nil {
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
			panic(err)
		}

		m, err := getMappingOrDefault(rt, key)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Location", m.Destination)
		if m.Permanent {
			w.WriteHeader(http.StatusPermanentRedirect)
		} else {
			w.WriteHeader(http.StatusTemporaryRedirect)
		}

		fmt.Fprintf(w, REDIRECT_BODY)
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
