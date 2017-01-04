package main

import (
	"bytes"
	"encoding/gob"
)

// MarshallBinary serializes an object into a slice of bytes.
func MarshallBinary(v interface{}) ([]byte, error) {
	w := &bytes.Buffer{}
	enc := gob.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

// UnmarshallBinary deserializes a slice of bytes into an object.
func UnmarshallBinary(b []byte, v interface{}) error {
	r := bytes.NewReader(b)
	dec := gob.NewDecoder(r)
	if err := dec.Decode(v); err != nil {
		return err
	}

	return nil
}
