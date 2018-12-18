package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/NYTimes/gziphandler"
	raven "github.com/getsentry/raven-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func LoadTemplateFuncs(name string, funcs template.FuncMap) *template.Template {
	tmpl := LoadTemplate(name)
	tmpl = tmpl.Funcs(funcs)
	return tmpl
}

func LoadTemplate(name string) *template.Template {
	tmpl := template.New("base.html")

	tmpl, err := tmpl.ParseFiles("templates/base.html", filepath.Join("templates", name+".html"))
	if err != nil {
		log.Fatalf("Couldn't load template '%s': %v", name, err)
	}

	return tmpl
}

func main() {
	portStr := os.Getenv("PORT")
	var port int
	if len(portStr) == 0 {
		port = 8080
	} else {
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Couldn't parse port %s", portStr)
		}
		port = portInt
	}

	r := mux.NewRouter()

	r.Methods("GET").Path("/").Handler(Entry())
	r.Methods("GET").Path("/about").Handler(About())
	r.Methods("POST").Path("/solve").Handler(Solve())

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	r.Use(raven.Recoverer)
	r.Use(gziphandler.GzipHandler)
	r.Use(func(h http.Handler) http.Handler { return handlers.LoggingHandler(os.Stdout, h) })

	log.Printf("Listening on localhost:%d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
