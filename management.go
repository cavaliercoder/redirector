package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type mgmtHandler struct {
	Runtime *Runtime
}

func ManagementHandler(rt *Runtime) http.Handler {
	return WrapHandler(rt, &mgmtHandler{rt})
}

func (c *mgmtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/stats/" {
		if r.Method != "GET" {
			panic(NewHTTPError(http.StatusMethodNotAllowed, nil))
		}

		c.getStatsHandler(w, r)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/mappings/") {
		switch r.Method {
		case "POST":
			c.postMappingHandler(w, r)
			return

		case "GET":
			c.getMappingsHandler(w, r)
			return

		case "DELETE":
			c.deleteMappingHandler(w, r)
			return

		default:
			panic(NewHTTPError(http.StatusMethodNotAllowed, nil))
		}
	}

	panic(NewHTTPError(http.StatusNotFound, nil))
}

func (c *mgmtHandler) getStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats := c.Runtime.Database.Stats()
	JSON(w, r, stats)
}

func (c *mgmtHandler) getMappingsHandler(w http.ResponseWriter, r *http.Request) {
	mappings, err := c.Runtime.Database.GetMappings()
	if err != nil {
		panic(err)
	}

	JSON(w, r, mappings)
}

func (c *mgmtHandler) postMappingHandler(w http.ResponseWriter, r *http.Request) {
	if m := r.Header.Get("Content-Type"); m != "application/json" {
		panic(NewHTTPError(http.StatusBadRequest, nil))
	}

	m := &Mapping{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(m); err != nil {
		panic(err)
	}

	if err := c.Runtime.Database.AddMapping(m); err != nil {
		panic(err)
	}

	w.Header().Set("Location", fmt.Sprintf("/mappings/%v", m.Key))
	w.WriteHeader(http.StatusCreated)
}

func (c *mgmtHandler) deleteMappingHandler(w http.ResponseWriter, r *http.Request) {
	m := Mapping{}
	fmt.Sscanf(r.URL.Path, "/mappings/%s", &m.Key)
	if m.Key == "" {
		if _, err := c.Runtime.Database.DeleteMappings(); err != nil {
			panic(err)
		}
	} else {
		if _, err := c.Runtime.Database.GetMapping(m.Key); err != nil {
			panic(err)
		}

		if err := c.Runtime.Database.DeleteMapping(m.Key); err != nil {
			panic(err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func serveManager(rt *Runtime) error {
	s := &http.Server{
		Addr:    rt.Config.MgmtAddr,
		Handler: ManagementHandler(rt),
	}

	rt.Logger.Printf("Listening for management commands on %v", rt.Config.MgmtAddr)
	return s.ListenAndServe()
}
