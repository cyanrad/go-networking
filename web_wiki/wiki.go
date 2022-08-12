package main

import (
	"log"
	"net/http"
)

func main() {
	//loading all template files instead of doing that
	//each time we do a request

	//tells http pkg to handle all requests to the web root ("/...") with handler
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createPageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil)) //listening to port 8080
	//blocks until prog is terminated
	//ListenAndServe will only return an err, so we log it if it happens
}
