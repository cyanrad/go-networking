package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const htmlFileName = "html/"

// Main handlers
//=====================
//view a wiki page. It will handle URLs prefixed with "/view/".
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Print(r)
	p, err := loadPage(title)
	if err != nil { //redirecting to edit page if does not exist
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		//adds an HTTP status code of http.StatusFound (302) and a Location header to the HTTP response.
		return
	}
	renderTemplate(w, "view", p)
}

//Edit page handler
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title) //loading page
	if err != nil {           //Error checking
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

//save request handler
//save request used in edit page, after clicking the save button [post]
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body") //the body from the textarea

	//converting to byte slice, from string
	p := &Page{Title: title, Body: template.HTML(body)}

	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//sends a specified HTTP response code (in this case "Internal Server Error")
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound) //redirecting to edit
}

type ExsistingPages struct {
	Pages []template.HTML
}

//index stuff
//=====================
//the handler for the index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(htmlFileName + "index.html") //reads content & returns for index
	if err != nil {                                            //err handling
		log.Fatal(err)
	}

	//creating a ExsistingPages type to send the wiki_pages to the
	execute_error := t.Execute(w, ExsistingPages{getExistingPages()})
	if execute_error != nil { //error handleing
		log.Fatal(execute_error)
	}
}

//the handelr for creating new pages from the create field in index.html
func createPageHandler(w http.ResponseWriter, r *http.Request) {
	pageTitle := r.FormValue("create_field")                  //the body from the text field (having the new page name)
	http.Redirect(w, r, "/edit/"+pageTitle, http.StatusFound) //redirects to the edit page as to create a new file at save
}
