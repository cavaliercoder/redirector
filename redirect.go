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

func Redirecthandler(rt *Runtime) http.Handler {
	return WrapHandler(rt, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, err := rt.Config.KeyBuilder.Parse(r)
		if err != nil {
			panic(err)
		}

		m, err := rt.Database.GetMapping(key)
		if err != nil {
			if err == MappingNotFoundError {
				err = NewHTTPError(http.StatusNotFound, err)
				panic(err)
			} else {
				panic(err)
			}
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Location", m.Destination)
		w.WriteHeader(http.StatusTemporaryRedirect)
		fmt.Fprintf(w, REDIRECT_BODY)
	}))
}

func serve(rt *Runtime) error {
	s := &http.Server{
		Addr:    rt.Config.ListenAddr,
		Handler: Redirecthandler(rt),
	}

	rt.Logger.Printf("Listening for redirect requests on %v", rt.Config.ListenAddr)
	return s.ListenAndServe()
}
