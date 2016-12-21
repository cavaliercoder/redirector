package main

type Mapping struct {
	Key         string `json:"key"`
	Destination string `json:"dest"`
	Permanent   bool   `json:"perm"`
}
