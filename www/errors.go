package main

import (
	"html/template"
	"net/http"
	"sync"
)

func WriteError(w http.ResponseWriter, status int, msg string, err error) {

	type HttpError struct {
		Status  int
		Message string
	}

	// load the templates
	var tmpl *template.Template
	var loadTmpl sync.Once
	loadTmpl.Do(func() {
		tmpl = template.New("error.html")
		tmpl = tmpl.Funcs(map[string]interface{}{
			"statusText": func(status int) string {
				return http.StatusText(status)
			},
		})
		tmpl, err = tmpl.ParseFiles("templates/error.html")
		if err != nil {
			panic(err)
		}
	})

	tmplErr := tmpl.Execute(w, HttpError{Status: status, Message: msg})
	if tmplErr != nil {
		panic(tmplErr)
	}
}

// Must panics if err != nil
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
