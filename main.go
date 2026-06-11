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
	"time"

	"github.com/gorilla/feeds"
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
	Content template.HTML
}

var Posts []blogData
var RSSFeed []byte

const siteURL = "https://dylanmeeks.engineer"

func parsePostDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Now()
	}
	return t
}

func generateRSS(posts []blogData) ([]byte, error) {
	now := time.Now()

	feed := &feeds.Feed{
		Title:       "Dylan's Blog",
		Link:        &feeds.Link{Href: siteURL + "/blog"},
		Description: "Posts from Dylan's personal website",
		Author:      &feeds.Author{Name: "Dylan Meeks"},
		Created:     now,
		Updated:     now,
	}

	for _, post := range posts {
		title := post.Tags["Title"]
		desc := post.Tags["Desc"]
		date := parsePostDate(post.Tags["Date"])

		url := siteURL + "/blog/" + title

		feed.Items = append(feed.Items, &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: url},
			Description: desc,
			Created:     date,
			Updated:     date,
			Content:     string(post.Content),
		})
	}

	gen, err := feed.ToRss()
	if err != nil {
		return nil, err
	}
	return []byte(gen), nil
}

func sortBlogs(blogs []blogData) []blogData {
	slices.SortFunc(blogs, func(a, b blogData) int {
		aTime := parsePostDate(a.Tags["Date"])
		bTime := parsePostDate(b.Tags["Date"])
		if aTime.Before(bTime) {
			return 1
		} else if bTime.Before(aTime) {
			return -1
		} else {
			return 0
		}
	})
	return blogs
}

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

		// start parsing custom tag fields
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

		// render rest of the file as the blog
		buf := make([]byte, 1024*516)
		_, err = reader.Read(buf)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read 'content' from file %s into buffer", e.Name())
		}
		contentHTML := MdToHTML([]byte(buf))
		blogs = append(blogs, blogData{tags, template.HTML(contentHTML)})
	}

	blogs = sortBlogs(blogs)

	// write CSS generated for blogs
	f, err := os.Create("static/css/chroma.css")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := WriteHighlightCSS(f); err != nil {
		panic(err)
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

	RSSFeed, err = generateRSS(Posts)
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
	http.HandleFunc("/blog/", blogPostHandler)
	http.HandleFunc("/rss.xml", rssHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func filterTags(tags map[string]string) map[string]string {
	newTags := maps.Clone(tags)
	for k := range newTags {
		if k == "Title" || k == "Date" || k == "Desc" {
			delete(newTags, k)
		}
	}
	return newTags
}

func getPost(posts []blogData, title string) blogData {
	for _, post := range posts {
		if title == post.Tags["Title"] {
			return post
		}
	}
	return blogData{}
}

func renderTemplate(w http.ResponseWriter, page string, data pageData) {
	files := []string{
		"templates/base.html",
		"templates/" + page,
	}

	tmpl, err := template.New("").Funcs(template.FuncMap{
		"filterTags": filterTags,
		"toLower":    strings.ToLower,
		"getPost":    getPost,
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

func blogPostHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/blog/")
	slug = strings.Trim(slug, "/")

	exists := false
	for _, post := range Posts {
		if slug == post.Tags["Title"] {
			exists = true
			break
		}
	}
	if !exists {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, "post.html", newPage("BLOG", slug, "2"))
}

func rssHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(RSSFeed)
}
