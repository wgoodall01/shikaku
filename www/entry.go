package main

import "net/http"

func Entry() http.HandlerFunc {
	tmpl := LoadTemplate("entry")

	return func(w http.ResponseWriter, r *http.Request) {
		Must(tmpl.Execute(w, nil))
	}
}
