package main

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"text/template"
)

// A Mapping maps a request key to a destination URL.
type Mapping struct {
	Key         string `json:"key"`
	Destination string `json:"dest"`
	Permanent   bool   `json:"perm,omitempty"`
	Comment     string `json:"comment,omitempty"`
	IsTemplate  bool   `json:"-"`
}

var (
	KeyMissingError             = fmt.Errorf("No key defined")
	DestinationMissingError     = fmt.Errorf("No destination defined")
	DestinationNotTemplateError = fmt.Errorf("Mapping destination is not a template")
)

var (
	mappingTemplates      = make(map[string]*template.Template)
	mappingTemplatesMutex = &sync.Mutex{}
)

func (m *Mapping) String() string {
	j := "->"
	if m.Permanent {
		j = "=>"
	}

	return fmt.Sprintf("%v %v %v", m.Key, j, m.Destination)
}

func (m *Mapping) Validate() error {
	if m.Key == "" {
		return KeyMissingError
	}

	if m.Destination == "" {
		return DestinationMissingError
	}

	if strings.Contains(m.Destination, "{{") {
		m.IsTemplate = true
	}

	return nil
}

// ComputeDestination expands any templated fields in the mapping destination
// URL.
//
// Request key may differ from actual mapping key when using 'catch-all'
// mappings such as the default mapping.
func (m *Mapping) ComputeDestination(vb ViewBag) (string, error) {
	if !m.IsTemplate {
		return "", DestinationNotTemplateError
	}

	// get template
	tmpl, err := func(m *Mapping) (*template.Template, error) {
		mappingTemplatesMutex.Lock()
		if tmpl, ok := mappingTemplates[m.Key]; ok {
			mappingTemplatesMutex.Unlock()
			return tmpl, nil
		}
		mappingTemplatesMutex.Unlock()

		if tmpl, err := template.New(m.Key).Parse(m.Destination); err != nil {
			return nil, err
		} else {
			mappingTemplatesMutex.Lock()
			mappingTemplates[m.Key] = tmpl
			mappingTemplatesMutex.Unlock()
			return tmpl, nil
		}
	}(m)
	if err != nil {
		return "", err
	}

	b := &bytes.Buffer{}
	if err := tmpl.Execute(b, vb); err != nil {
		return "", err
	}

	return string(b.Bytes()), nil
}
