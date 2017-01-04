package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// returns a http.Client that does not follow redirects
func testHttpClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func testRedirectServer(fn func(*Runtime, *httptest.Server)) {
	tmpBoltDB(func(db Database) {
		rt := &Runtime{
			Config:   &Config{},
			Database: db,
			Logger:   log.New(ioutil.Discard, "", 0),
		}

		ts := httptest.NewServer(RedirectHandler(rt))
		defer ts.Close()

		fn(rt, ts)
	})
}

func TestDefaultKey(t *testing.T) {
	testRedirectServer(func(rt *Runtime, ts *httptest.Server) {
		//rt.Config.DefaultKey = "default"
		rt.Config.KeyBuilder = RequestURIPathKeyBuilder()

		// add default mapping
		m := &Mapping{
			Key:         "default",
			Destination: "/okay",
			Permanent:   true,
		}
		if err := rt.Database.AddMapping(m); err != nil {
			panic(err)
		}

		res, err := testHttpClient().Get(ts.URL + "/does/not/exist")
		if err != nil {
			panic(err)
		}

		if res.StatusCode != http.StatusMovedPermanently {
			t.Fatalf("Expected default mapping with status %v, got %v", http.StatusMovedPermanently, res.StatusCode)
		}

		loc := res.Header.Get("Location")
		if loc != m.Destination {
			t.Fatalf("Expected default mapping to '%v', got '%v'", m.Destination, loc)
		}
	})
}
