package main

import (
	"html/template"
	"net/http"
	"regexp"
)

//Used so that we don't reload the template each time we do a request
//Must() panics at non-nil err value
var templates = template.Must(template.ParseFiles("html/edit.html", "html/view.html"))

//Generating webpage html from .html template files
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p) //Generates actual HTML page
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// writing the generated HTML to the http.ResponseWriter.
	// The .Title and .Body dotted identifiers refer to p.Title and p.Body
}

//A regexp to check for a valid path/link
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

//creats the handlers to reduce redundency
//checking the path correctness ->
//running the handler function
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//If the title is valid, it will be returned along with a nil error value.
		//else write a "404 Not Found" error to the HTTP connection
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
