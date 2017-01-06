package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
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
		for _, m := range testMappings {
			if err := db.AddMapping(&m); err != nil {
				panic(err)
			}
		}

		rt := &Runtime{
			Config: &Config{
				ExitOnError: true,
				KeyBuilder:  RequestURIPathKeyBuilder(),
			},
			Database:     db,
			Logger:       log.New(ioutil.Discard, "", 0),
			AccessLogger: log.New(ioutil.Discard, "", 0),
		}
		InitTemplates()

		ts := httptest.NewServer(RedirectHandler(rt))
		defer ts.Close()

		fn(rt, ts)
	})
}

func TestMissingKey(t *testing.T) {
	testRedirectServer(func(rt *Runtime, ts *httptest.Server) {
		res, err := testHttpClient().Get(ts.URL + "/does/not/exist")
		if err != nil {
			panic(err)
		}

		if res.StatusCode != http.StatusNotFound {
			t.Fatalf("Expected missing mapping with status %v, got %v", http.StatusNotFound, res.StatusCode)
		}
	})
}

func TestPermanentKey(t *testing.T) {
	testRedirectServer(func(rt *Runtime, ts *httptest.Server) {
		dest := "/okay"
		res, err := testHttpClient().Get(ts.URL + "/permanent")
		if err != nil {
			panic(err)
		}

		if res.StatusCode != http.StatusPermanentRedirect {
			t.Fatalf("Expected permanent mapping with status %v, got %v", http.StatusPermanentRedirect, res.StatusCode)
		}

		loc := res.Header.Get("Location")
		if loc != dest {
			t.Fatalf("Expected permanent mapping to '%v', got '%v'", dest, loc)
		}
	})
}

func TestTemporaryKey(t *testing.T) {
	testRedirectServer(func(rt *Runtime, ts *httptest.Server) {
		dest := "/okay"
		res, err := testHttpClient().Get(ts.URL + "/temporary")
		if err != nil {
			panic(err)
		}

		if res.StatusCode != http.StatusTemporaryRedirect {
			t.Fatalf("Expected temporary mapping with status %v, got %v", http.StatusTemporaryRedirect, res.StatusCode)
		}

		loc := res.Header.Get("Location")
		if loc != dest {
			t.Fatalf("Expected temporary mapping to '%v', got '%v'", dest, loc)
		}
	})
}

func TestDefaultKey(t *testing.T) {
	testRedirectServer(func(rt *Runtime, ts *httptest.Server) {
		dest := "/okay"
		rt.Config.DefaultKey = "default"
		rt.Config.KeyBuilder = RequestURIPathKeyBuilder()

		// test non-existant mapping
		res, err := testHttpClient().Get(ts.URL + "/does/not/exist")
		if err != nil {
			panic(err)
		}

		if res.StatusCode != http.StatusTemporaryRedirect {
			t.Fatalf("Expected default mapping with status %v, got %v", http.StatusTemporaryRedirect, res.StatusCode)
		}

		loc := res.Header.Get("Location")
		if loc != dest {
			t.Fatalf("Expected default mapping to '%v', got '%v'", dest, loc)
		}

		// test existing mapping still works
		for _, m := range testMappings {
			u, _ := url.Parse(ts.URL)
			u.Path = filepath.Join(u.Path, m.Key)
			res, err = testHttpClient().Get(u.String())
			if err != nil {
				panic(err)
			}

			loc = res.Header.Get("Location")
			if loc != m.Destination {
				t.Fatalf("Bad mapping destination '%v', expected '%v'", loc, m.Destination)
			}
		}
	})
}
