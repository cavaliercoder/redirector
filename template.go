package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

const (
	defaultTemplate = `<html>
<head><title>{{.status}} {{.statusText}}</title></head>
<body bgcolor="white">
<center><h1>{{.status}} {{.statusText}}</h1></center>
<hr><center>` + PACKAGE_NAME + `/` + PACKAGE_VERSION + `</center>
</body>
</html>`
)

var BodyNotFoundError = fmt.Errorf("No body content found for the given status code")

var (
	statusCodes = []int{
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,
		http.StatusNotFound,
		http.StatusInternalServerError,
	}

	bodies = make(map[int]string, 0)
)

func InitTemplates() error {
	tmpl, err := template.New("default").Parse(defaultTemplate)
	if err != nil {
		return err
	}

	for _, code := range statusCodes {
		b := &bytes.Buffer{}
		data := map[string]string{
			"status":     fmt.Sprintf("%v", code),
			"statusText": http.StatusText(code),
		}

		if err := tmpl.Execute(b, data); err != nil {
			return err
		}

		bodies[code] = b.String()
	}

	return nil
}

func BodyForStatus(code int) (string, error) {
	if body, ok := bodies[code]; ok {
		return body, nil
	}

	panic(code)

	return "", BodyNotFoundError
}
