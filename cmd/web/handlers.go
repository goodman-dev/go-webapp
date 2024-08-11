package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

// Our first handler, used to serve the root
func home(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/base.html.tmpl",
		"./ui/html/pages/home.html.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

}

func snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snipper with ID %d...", id)
}

func snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating snippets"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Creating new snippet"))
}
