package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/wgoodall01/shikaku"
)

type solveHandler struct {
	http.Handler

	tmpl *template.Template
}

func Solve() http.Handler {
	tmpl := LoadTemplateFuncs("solve", map[string]interface{}{
		"add": func(a, b int) int {
			return a + b
		},
	})

	return &solveHandler{
		tmpl: tmpl,
	}
}

func (h *solveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		WriteError(w, 400, "Couldn't parse form", err)
		return
	}

	if r.Form["rows"] == nil || r.Form["cols"] == nil {
		WriteError(w, 400, "Invalid request", nil)
		return
	}

	rows, err := strconv.Atoi(r.Form["rows"][0])
	if err != nil {
		WriteError(w, 400, "Invalid number of rows", err)
		return
	}

	cols, err := strconv.Atoi(r.Form["cols"][0])
	if err != nil {
		WriteError(w, 400, "Invalid num. cols", err)
		return
	}

	if cols > 40 || rows > 40 {
		WriteError(w, 400, fmt.Sprintf("%d by %d board is way too large.", rows, cols), nil)
		return
	}

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
			if err != nil {
				WriteError(w, 400, fmt.Sprintf("Bad square at [%d,%d]", c, r), err)
				return
			}

			if val > rows*cols {
				WriteError(w, 400, fmt.Sprintf("Square [%d,%d] is too large", c, r), nil)
				return
			}

			bo.Grid[r][c] = shikaku.NewGiven(val)
		}
	}

	// Solve the puzzle
	tStart := time.Now()
	solveErr := bo.Solve()
	duration := time.Since(tStart).Seconds() / 1000 // in ms

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
		Err      error
		Soln     template.HTML
		Small    bool
		Duration float64
	}{
		Err:      solveErr,
		Soln:     template.HTML(buf.String()),
		Small:    rows > 15 || cols > 15,
		Duration: duration,
	}

	if err := h.tmpl.Execute(w, viewState); err != nil {
		panic(err)
	}
}
