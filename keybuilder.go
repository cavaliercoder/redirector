package main

import (
	"fmt"
	"net/http"
)

var ParamNotFoundError = fmt.Errorf("Key parameter not found in request URI")

// A KeyBuilder translates a client HTTP request into a URL mapping key that may
// be used to lookup the destination URL of a redirect mapping.
type KeyBuilder interface {
	Parse(r *http.Request) (string, error)
}

// The KeyBuilderFunc type is an adapter to allow the use of ordinary functions
// as KeyBuilders. If f is a function with the appropriate signature,
// KeyBuilderFunc(f) is a KeyBuilder that calls f.
type KeyBuilderFunc func(r *http.Request) (string, error)

// Parse calls f(r).
func (f KeyBuilderFunc) Parse(r *http.Request) (string, error) {
	return f(r)
}

// RequestURIKeyBuilder returns a KeyBuilder that uses the full request URI as
// a mapping key.
func RequestURIKeyBuilder() KeyBuilder {
	return KeyBuilderFunc(func(r *http.Request) (string, error) {
		if r.URL.Scheme == "" {
			if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
				r.URL.Scheme = scheme
			} else {
				r.URL.Scheme = "http"
			}
		}
		if r.URL.Host == "" {
			if host := r.Header.Get("Host"); host != "" {
				r.URL.Host = host
			} else if host := r.Header.Get("X-Forwarded-Host"); host != "" {
				r.URL.Host = host
			} else {
				r.URL.Host = "."
			}
		}
		return r.URL.String(), nil
	})
}

// RequestURIPathKeyBuilder returns a KeyBuilder that uses the request URI path
// component (e.g. "/some/path") as a mapping key.
func RequestURIPathKeyBuilder() KeyBuilder {
	return KeyBuilderFunc(func(r *http.Request) (string, error) {
		return r.URL.Path, nil
	})
}

// RequestParamKeyBuilder returns a KeyBuilder that uses a given request URI
// query parameter (e.g. "?key=some_key") as a mapping key.
func RequestParamKeyBuilder(param string) KeyBuilder {
	return KeyBuilderFunc(func(r *http.Request) (string, error) {
		key := r.URL.Query().Get(param)
		if key == "" {
			return "", NewHTTPError(http.StatusNotFound, ParamNotFoundError)
		}

		return key, nil
	})
}
