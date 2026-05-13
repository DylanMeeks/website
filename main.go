package main

import (
	// "errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	// "os"
	// "regexp"
)

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/skills", skillsHandler)
	http.HandleFunc("/about", aboutHandler)

	fmt.Println("Server running on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.html")
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "projects.html")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contact.html")
}

func skillsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "skills.html")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.html")
}
