package main

import (
	"fmt"
)

// A Mapping maps a request key to a destination URL.
type Mapping struct {
	Key         string `json:"key"`
	Destination string `json:"dest"`
	Permanent   bool   `json:"perm,omitempty"`
	Comment     string `json:"comment,omitempty"`
}

func (m *Mapping) String() string {
	j := "->"
	if m.Permanent {
		j = "=>"
	}

	return fmt.Sprintf("%v %v %v", m.Key, j, m.Destination)
}
