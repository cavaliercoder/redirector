package main

import (
	"net/http"
)

type KeyBuilder interface {
	Parse(r *http.Request) (string, error)
}

type KeyBuilderFunc func(r *http.Request) (string, error)

func (f KeyBuilderFunc) Parse(r *http.Request) (string, error) {
	return f(r)
}

func RequestURIKeyBuilder() KeyBuilder {
	return KeyBuilderFunc(func(r *http.Request) (string, error) {
		return r.URL.String(), nil
	})
}

func RequestURIPathKeyBuilder() KeyBuilder {
	return KeyBuilderFunc(func(r *http.Request) (string, error) {
		return r.URL.Path, nil
	})
}

func RequestParamKeyBuilder(param string) KeyBuilder {
	return KeyBuilderFunc(func(r *http.Request) (string, error) {
		key := r.URL.Query().Get(param)
		if key == "" {
			return "", NewHTTPError(http.StatusNotFound, nil)
		}

		return key, nil
	})
}
