package main

import "net/http"

func About() http.HandlerFunc {
	tmpl := LoadTemplate("about")
	return func(w http.ResponseWriter, r *http.Request) {
		Must(tmpl.Execute(w, nil))
	}
}
