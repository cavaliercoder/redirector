package main

import (
	"net/http"
)

const (
	REDIRECT_BODY = `<html>
<head><title>301 Moved Permanently</title></head>
<body bgcolor="white">
<center><h1>301 Moved Permanently</h1></center>
<hr><center>` + PACKAGE_NAME + `/` + PACKAGE_VERSION + `</center>
</body>
</html>`
)

type RedirectHandler struct {
	Runtime *Runtime
}

func (c *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key, err := c.Runtime.Config.KeyBuilder.Parse(r)
	if err != nil {
		panic(err)
	}

	m, err := c.Runtime.Database.GetMapping(key)
	if err != nil {
		if err == MappingNotFoundError {
			http.NotFound(w, r)
			return
		} else {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Location", m.Destination)
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(REDIRECT_BODY))
}

func serve(rt *Runtime) error {
	srv := &RedirectHandler{rt}

	s := &http.Server{
		Addr:    rt.Config.ListenAddr,
		Handler: srv,
	}

	rt.Logger.Printf("Listening for redirect requests on %v", rt.Config.ListenAddr)
	return s.ListenAndServe()
}
