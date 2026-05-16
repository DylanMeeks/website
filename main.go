package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {

	http.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))),
	)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/skills", skillsHandler)
	http.HandleFunc("/about", aboutHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.ParseFiles(
		"templates/base.html",
		"templates/"+tmpl,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "home.html", map[string]string{
		"Title": "Home",
	})
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "projects.html", map[string]string{
		"Title": "Projects",
	})
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contact.html", map[string]string{
		"Title": "Contact",
	})
}

func skillsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "skills.html", map[string]string{
		"Title": "Skills",
	})
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.html", map[string]string{
		"Title": "About",
	})
}
