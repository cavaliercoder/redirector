package main

import (
	"fmt"
)

type Mapping struct {
	Key         string `json:"key"`
	Destination string `json:"dest"`
	Permanent   bool   `json:"perm"`
}

func (m *Mapping) String() string {
	j := "->"
	if m.Permanent {
		j = "=>"
	}

	return fmt.Sprintf("%v %v %v", m.Key, j, m.Destination)
}
