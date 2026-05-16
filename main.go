package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type pageData struct {
	Title      string
	ManName    string
	ManSection string
	ManVolume  string
	ManSource  string
	ManDate    string
}

func newPage(name, title string) pageData {
	return pageData{
		Title:      title,
		ManName:    name,
		ManSection: "1",
		ManVolume:  "Dylan's Personal Manual",
		ManSource:  "dylan 1.0",
		ManDate:    "May 2026",
	}
}

func main() {
	http.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))),
	)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/skills", skillsHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/contact", contactHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderTemplate(w http.ResponseWriter, page string, data pageData) {
	t, err := template.ParseFiles(
		"templates/base.html",
		"templates/"+page,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "home.html", newPage("DYLAN", "dylan"))
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "projects.html", newPage("PROJECTS", "projects"))
}

func skillsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "skills.html", newPage("SKILLS", "skills"))
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.html", newPage("ABOUT", "about"))
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contact.html", newPage("CONTACT", "contact"))
}
