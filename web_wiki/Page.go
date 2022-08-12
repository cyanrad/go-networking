package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
)

//A wiki is consists of interconnected pages,
//Each of which has a title and a body
type Page struct {
	Title string
	Body  template.HTML //byte slice
	// we use byte as it is the type expected by the io library
}

//presistent storage (save body to text file, with title as file name)
func (p *Page) save() error {
	filename := "wiki_pages/" + p.Title + ".txt" //physical file name
	fmt.Println("Saving Page:", filename)
	//writing bytes to a file
	return os.WriteFile(filename, []byte(p.Body), 0600) //write file returns an error type (nil if all goes well)
	//0600 means that it should be created with read-write perms only for the user.
}

//loading page from file
func loadPage(title string) (*Page, error) {
	filename := "wiki_pages/" + title + ".txt"
	fmt.Println("Loading Page:", filename)
	body, err := os.ReadFile(filename)

	if err != nil { //if file has error/ doesn't exist
		return nil, err
	}
	return &Page{title, template.HTML(body)}, nil //creating a new page object
}

//checks all the existing files in the wiki_pages dir
func getExistingPages() []template.HTML {
	files, err := os.ReadDir("wiki_pages")
	if err != nil {
		log.Fatal(err)
	}

	var returnStrSlice []template.HTML
	for _, file := range files {
		tempFileName := file.Name()[:len(file.Name())-4]
		returnStrSlice = append(returnStrSlice, template.HTML("<li><a href=\"/view/"+tempFileName+"\">"+tempFileName+"</a></li>"))
	}
	return returnStrSlice
}
