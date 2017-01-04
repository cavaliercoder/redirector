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
	mappings := []Mapping{
		{"default", "/okay", true},
		{"/test1", "/okay-test1", true},
		{"/test2", "/okay-test2", true},
	}

	tmpBoltDB(func(db Database) {
		for _, m := range mappings {
			if err := db.AddMapping(&m); err != nil {
				panic(err)
			}
		}

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
		defaultDest := "/okay"
		rt.Config.DefaultKey = "default"
		rt.Config.KeyBuilder = RequestURIPathKeyBuilder()

		// test non-existant mapping
		res, err := testHttpClient().Get(ts.URL + "/does/not/exist")
		if err != nil {
			panic(err)
		}

		if res.StatusCode != http.StatusMovedPermanently {
			t.Fatalf("Expected default mapping with status %v, got %v", http.StatusMovedPermanently, res.StatusCode)
		}

		loc := res.Header.Get("Location")
		if loc != defaultDest {
			t.Fatalf("Expected default mapping to '%v', got '%v'", defaultDest, loc)
		}

		// test existing mapping still works
		res, err = testHttpClient().Get(ts.URL + "/test1")
		if err != nil {
			panic(err)
		}

		if res.StatusCode != http.StatusMovedPermanently {
			t.Fatalf("Expected real mapping with status %v, got %v", http.StatusMovedPermanently, res.StatusCode)
		}

		loc = res.Header.Get("Location")
		if loc != "/okay-test1" {
			t.Fatalf("Expected real mapping to '%v', got '%v'", "/okay-test1", loc)
		}
	})
}
