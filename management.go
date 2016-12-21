package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ManagementHandler struct {
	Runtime *Runtime
}

func (c *ManagementHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger := c.Runtime.Logger
	status := http.StatusInternalServerError

	defer func() {
		if r.Body != nil {
			r.Body.Close()
		}

		duration := time.Since(start)
		/*
		 * if err := recover(); err != nil {
		 *    logger.Printf("%v", err)
		 * }
		 */

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(status)
		fmt.Fprintf(w, "%d %s\n", status, http.StatusText(status))

		logger.Printf("%s %s %d %v\n", r.Method, r.URL.Path, status, duration)
	}()

	switch r.URL.Path {
	case "/stats/":
		if r.Method == "GET" {
			status = c.getStatsHandler(w, r)
		} else {
			status = http.StatusMethodNotAllowed
		}

	case "/mappings/":
		switch r.Method {
		case "POST":
			status = c.postMappingHandler(w, r)

		case "GET":
			status = c.getMappingsHandler(w, r)

		default:
			status = http.StatusMethodNotAllowed
		}

	default:
		status = http.StatusNotFound
	}
}

func (c *ManagementHandler) JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		panic(err)
	}
}

func (c *ManagementHandler) getStatsHandler(w http.ResponseWriter, r *http.Request) int {
	stats := c.Runtime.Database.Stats()
	c.JSON(w, r, stats)

	return http.StatusOK
}

func (c *ManagementHandler) getMappingsHandler(w http.ResponseWriter, r *http.Request) int {
	mappings, err := c.Runtime.Database.GetMappings()
	if err != nil {
		panic(err)
	}

	c.JSON(w, r, mappings)
	return http.StatusOK
}

func (c *ManagementHandler) postMappingHandler(w http.ResponseWriter, r *http.Request) int {
	if m := r.Header.Get("Content-Type"); m != "application/json" {
		return http.StatusBadRequest
	}

	m := &Mapping{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(m); err != nil {
		panic(err)
	}

	if err := c.Runtime.Database.AddMapping(m); err != nil {
		panic(err)
	}

	return http.StatusCreated
}

func serveManager(rt *Runtime) error {
	srv := &ManagementHandler{rt}

	s := &http.Server{
		Addr:    rt.Config.MgmtAddr,
		Handler: srv,
	}

	rt.Logger.Printf("Listening for management commands on %v", rt.Config.MgmtAddr)
	return s.ListenAndServe()
}
