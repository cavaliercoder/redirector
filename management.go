package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	startTime time.Time
)

type RuntimeStats struct {
	Status   string        `json:"status"`
	Version  string        `json:"version"`
	Uptime   int64         `json:"uptime"`
	Database DatabaseStats `json:"database"`
}

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

	if r.URL.Path == "/config/" {
		if r.Method != "GET" {
			panic(NewHTTPError(http.StatusMethodNotAllowed, nil))
		}

		JSON(w, r, c.Runtime.Config)
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
	dbstats, err := c.Runtime.Database.Stats()
	if err != nil {
		panic(err)
	}

	stats := &RuntimeStats{
		Status:   "OK",
		Version:  PACKAGE_VERSION,
		Uptime:   int64(time.Since(startTime).Seconds()),
		Database: dbstats,
	}

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

	mappings := make([]*Mapping, 0)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&mappings); err != nil {
		panic(err)
	}

	for i, m := range mappings {
		if err := c.Runtime.Database.AddMapping(m); err != nil {
			panic(fmt.Errorf("Error adding mapping '%v at index [%v]': %v", m.Key, i, err))
		}
	}

	c.Runtime.Logger.Printf("Added %v mappings", len(mappings))
	if len(mappings) == 1 {
		w.Header().Set("Location", fmt.Sprintf("/mappings/%v", mappings[0].Key))
	}

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
	startTime = time.Now()
	s := &http.Server{
		Addr:    rt.Config.MgmtAddr,
		Handler: ManagementHandler(rt),
	}

	rt.Logger.Printf("Listening for management commands on %v", rt.Config.MgmtAddr)
	return s.ListenAndServe()
}
