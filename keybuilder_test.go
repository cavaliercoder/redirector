package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestPathKeyBuilder(t *testing.T) {
	expect := "/test/path/"
	c := &Config{
		KeyBuilderName: "path",
	}

	if err := c.initialize(); err != nil {
		panic(err)
	}

	if c.KeyBuilder == nil {
		t.Fatalf("KeyBuilder not instanciated")
	}

	r, err := http.NewRequest("GET", fmt.Sprintf("http://path.test%v", expect), nil)
	if err != nil {
		panic(err)
	}

	key, err := c.KeyBuilder.Parse(r)
	if err != nil {
		panic(err)
	}

	if key != expect {
		t.Fatalf("Expected key %v, got %v", expect, key)
	}
}

func TestURIKeyBuilder(t *testing.T) {
	expect := "http://test.local/some/path/?query=value&query2=value2"
	kb := RequestURIKeyBuilder()

	t.Run("Direct", func(t *testing.T) {
		req, _ := http.NewRequest("GET", expect, nil)
		key, err := kb.Parse(req)
		if err != nil {
			panic(err)
		}
		if key != expect {
			t.Errorf("expected key: %v, got: %v", expect, key)
		}
	})

	t.Run("VHost", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/some/path/?query=value&query2=value2", nil)
		req.Header.Set("host", "test.local")
		key, err := kb.Parse(req)
		if err != nil {
			panic(err)
		}
		if key != expect {
			t.Errorf("expected key: %v, got: %v", expect, key)
		}
	})
}

func TestParamKeyBuilderConfig(t *testing.T) {
	expect := "/test"
	c := &Config{
		KeyBuilderName: "param:key",
	}

	if err := c.initialize(); err != nil {
		panic(err)
	}

	if c.KeyBuilder == nil {
		t.Fatalf("KeyBuilder not instanciated")
	}

	r, err := http.NewRequest("GET", fmt.Sprintf("?key=%v&not-key=/whatevs", expect), nil)
	if err != nil {
		panic(err)
	}

	key, err := c.KeyBuilder.Parse(r)
	if err != nil {
		panic(err)
	}

	if key != expect {
		t.Fatalf("Expected key %v, got %v", expect, key)
	}
}

func TestParamKeyBuilder(t *testing.T) {
	expect := "/test"
	kb := RequestParamKeyBuilder("key")

	r, err := http.NewRequest("GET", fmt.Sprintf("?key=%v", expect), nil)
	if err != nil {
		panic(err)
	}

	key, err := kb.Parse(r)
	if err != nil {
		panic(err)
	}

	if key != expect {
		t.Fatalf("Expected key %v, got %v", expect, key)
	}
}

func TestMissingParamKeyBuilder(t *testing.T) {
	kb := RequestParamKeyBuilder("key")

	r, err := http.NewRequest("GET", "http://no-param.test/", nil)
	if err != nil {
		panic(err)
	}

	key, err := kb.Parse(r)
	if herr, ok := err.(*HTTPError); ok {
		if herr.Err != ParamNotFoundError {
			t.Fatalf("Expected error '%v', got '%v'", ParamNotFoundError, herr.Err)
		}
	} else {
		t.Fatalf("Expected HTTP Not Found error, got '%v'", err)
	}

	if key != "" {
		t.Fatalf("Expected empty key, got '%v'", key)
	}
}
