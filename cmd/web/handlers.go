package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/goodman-dev/go-webapp/internal/models"
)

// Our first handler, used to serve the root
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// w.Header().Add("Server", "Go")

	// files := []string{
	// 	"./ui/html/base.html.tmpl",
	// 	"./ui/html/pages/home.html.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// }

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// %v prints in the default format when printing structs,
	// the + adds field names https://pkg.go.dev/fmt
	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating snippets"))
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	title := "0 smail"
	content := "0 snail\nClimb Mount Fuji,\n"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
