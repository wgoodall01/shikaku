package main

import (
	"bytes"
	"fmt"
	"github.com/wgoodall01/shikaku"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type HttpError struct {
	Status  int
	Message string
	More    []string
}

func ErrStatus(status int) HttpError {
	return HttpError{
		Status:  status,
		Message: http.StatusText(status),
		More:    []string{},
	}
}

func ErrMsg(status int, msg string) HttpError {
	return HttpError{
		Status:  status,
		Message: http.StatusText(status),
		More:    []string{msg},
	}
}

// Must panics if err != nil, with an internal HttpError
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// MustErr panics if err != nil, with a custom HttpError throw
func MustErr(err error, throw HttpError) {
	if err != nil {
		panic(throw)
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(status int) {
	sw.status = status
	sw.ResponseWriter.WriteHeader(status)
}

func (sw *statusWriter) Write(buf []byte) (int, error) {
	sw.status = 200
	return sw.ResponseWriter.Write(buf)
}

func WrapMux(h http.Handler) http.Handler {
	errorTmpl := template.Must(template.ParseFiles("templates/error.html"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{
			ResponseWriter: w,
		}

		defer func(start time.Time) {
			err := recover()

			if err != nil {
				var he HttpError
				switch v := err.(type) {
				case HttpError:
					he = v
				default:
					if sw.status != 0 {
						he = ErrStatus(sw.status)
					} else {
						he = ErrStatus(500)
					}
				}

				if sw.status != he.Status && sw.status == 0 {
					sw.WriteHeader(he.Status)
				} else {
					sw.status = he.Status
				}

				tmpErr := errorTmpl.Execute(sw, he)
				if tmpErr != nil {
					sw.Write([]byte("Error rendering error. Send help."))
				}
			}

			log.Printf(
				"%v\t%s\t%-30s\t%d %.2fms",
				r.RemoteAddr,
				r.Method,
				r.URL,
				sw.status,
				time.Since(start).Seconds()*1000,
			)

			if err != nil {
				log.Printf("\terror:%v", err)
			}

		}(time.Now())

		h.ServeHTTP(sw, r)
	})
}

func LoadTemplateFuncs(name string, funcs template.FuncMap) *template.Template {
	tmpl := template.New("base.html")

	tmpl = tmpl.Funcs(funcs)

	tmpl, err := tmpl.ParseFiles("templates/base.html", filepath.Join("templates", name+".html"))
	if err != nil {
		log.Fatalf("Couldn't load template '%s': %v", name, err)
	}

	return tmpl
}

func LoadTemplate(name string) *template.Template {
	return LoadTemplateFuncs(name, nil)
}

func Entry() http.HandlerFunc {
	tmpl := LoadTemplate("entry")

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			panic(ErrStatus(404))
			return
		}

		switch r.Method {
		case "GET":
			Must(tmpl.Execute(w, nil))
		default:
			panic(ErrStatus(404))
		}

	}
}

func Solve() http.HandlerFunc {
	tmpl := LoadTemplateFuncs("solve", map[string]interface{}{
		"add": func(a, b int) int {
			return a + b
		},
	})

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			MustErr(r.ParseForm(), ErrMsg(400, "Couldn't parse form"))

			rows, err := strconv.Atoi(r.Form["rows"][0])
			MustErr(err, ErrMsg(400, "Invalid num. rows"))

			cols, err := strconv.Atoi(r.Form["cols"][0])
			MustErr(err, ErrMsg(400, "Invalid num. cols"))

			// Allocate board
			bo := &shikaku.Board{}
			for r := 0; r < rows; r++ {
				bo.Grid = append(bo.Grid, make([]shikaku.Square, cols))
			}

			for i, valStr := range r.Form["e"] {
				r := int(i / cols)
				c := i % cols
				if len(valStr) == 0 {
					bo.Grid[r][c] = shikaku.NewBlank()
				} else {
					val, err := strconv.Atoi(valStr)
					MustErr(err, ErrMsg(400, fmt.Sprintf("Bad cell at [%d,%d]", c, r)))
					bo.Grid[r][c] = shikaku.NewGiven(val)
				}
			}

			// Solve the puzzle
			solveErr := bo.Solve()

			// Build the table.
			buf := bytes.Buffer{}
			fmt.Fprint(&buf, "<thead>")
			fmt.Fprint(&buf, "<td></td>")
			for i := 0; i < bo.Width(); i++ {
				fmt.Fprintf(&buf, `<td class="solve_label">%d</td>`, i+1)
			}
			fmt.Fprint(&buf, "</thead>")

			var pos shikaku.Vec2
			fmt.Fprint(&buf, "<tbody>")
			for pos[1] = 0; pos[1] < bo.Height(); pos[1]++ {
				fmt.Fprint(&buf, "<tr>")
				fmt.Fprintf(&buf, `<td class="solve_label">%d</td>`, pos[1]+1)
				for pos[0] = 0; pos[0] < bo.Width(); pos[0]++ {
					sq := bo.Get(pos)
					if shikaku.IsNotFinal(*sq) {
						// Write empty square
						fmt.Fprint(&buf, `<td class="solve_empty"></td>`)
					} else if sq.Final.A == pos {
						// It's the top-left, write a cell w/ colspan and rowspan.
						colspan := sq.Final.Width()
						rowspan := sq.Final.Height()
						fmt.Fprintf(
							&buf,
							`<td class="solve_rect" colspan="%d" rowspan="%d">%d</td>`,
							colspan,
							rowspan,
							bo.Get(sq.Final.Given).Area,
						)
					}
				}
				fmt.Fprint(&buf, "</tr>")
			}
			fmt.Fprintf(&buf, "</tbody>")

			viewState := struct {
				Err  error
				Soln template.HTML
			}{
				Err:  solveErr,
				Soln: template.HTML(buf.String()),
			}

			Must(tmpl.Execute(w, viewState))
		default:
			http.Redirect(w, r, "/", 303)
		}
	}
}

func About() http.HandlerFunc {
	tmpl := LoadTemplate("about")
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			Must(tmpl.Execute(w, nil))
		default:
			panic(ErrStatus(404))
		}
	}
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

	http.HandleFunc("/", Entry())
	http.HandleFunc("/solve", Solve())
	http.HandleFunc("/about", About())

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Printf("Listening on localhost:%d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), WrapMux(http.DefaultServeMux))
}
