package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type pageData struct {
	Title      string
	ManName    string
	ManSection string
	ManVolume  string
	ManSource  string
	ManDate    string
	Posts      []blogData
}

func newPage(name, title, section string) pageData {
	return pageData{
		Title:      title,
		ManName:    name,
		ManVolume:  "Dylan's Personal Manual",
		ManSource:  "dylan 1.0",
		ManSection: "1",
		ManDate:    "May 2026",
		Posts:      Posts,
	}
}

type blogData struct {
	Tags    map[string]string
	Content string
}

func filterTags(tags map[string]string) map[string]string {
	newTags := maps.Clone(tags)
	for k, _ := range newTags {
		if k == "Title" || k == "Date" || k == "Desc" {
			delete(newTags, k)
		}
	}
	return newTags
}

var Posts []blogData

func genBlogs() ([]blogData, error) {
	dir, err := os.ReadDir("content")
	if err != nil {
		return nil, fmt.Errorf("failed to open content dir")
	}
	blogs := make([]blogData, 0)
	for _, e := range dir {
		if e.Type().IsDir() {
			continue
		}
		f, err := os.Open(filepath.Join("content", e.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s", e.Name())
		}
		defer f.Close()

		reader := bufio.NewReader(f)
		start, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read starting '---' file %s", e.Name())
		}
		if !slices.Equal(start, []byte{'-', '-', '-', '\n'}) {
			return nil, fmt.Errorf("failed to find '---' at top of file %s", e.Name())
		}
		tags := make(map[string]string)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				return nil, fmt.Errorf("failed to read line in file %s", e.Name())
			}
			if slices.Equal(line, []byte{'-', '-', '-', '\n'}) {
				break
			}
			strs := strings.Split(string(line), ":")
			tags[strings.Trim(strs[0], " \n")] = strings.Trim(strs[1], " \n")
		}

		buf := make([]byte, 1024*516)
		_, err = reader.Read(buf)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read 'content' from file %s into buffer", e.Name())
		}
		contentHTML := MdToHTML([]byte(buf))
		blogs = append(blogs, blogData{tags, string(contentHTML)})
	}
	return blogs, nil
}

func main() {
	InitRenderer()
	var err error
	Posts, err = genBlogs()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))),
	)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/skills", skillsHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/blog", blogHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderTemplate(w http.ResponseWriter, page string, data pageData) {
	files := []string{
		"templates/base.html",
		"templates/" + page,
	}

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"filterTags": filterTags,
		"toLower":    strings.ToLower,
	}).ParseFiles(files...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "home.html", newPage("DYLAN", "dylan", "1"))
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "projects.html", newPage("PROJECTS", "projects", "1"))
}

func skillsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "skills.html", newPage("SKILLS", "skills", "1"))
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.html", newPage("ABOUT", "about", "1"))
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contact.html", newPage("CONTACT", "contact", "1"))
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "blog.html", newPage("BLOG", "blog", "1"))
}
