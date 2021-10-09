package main

import (
	"html/template"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var links = map[string]string{
	"bitstamp": "https://www.bitstamp.net/",
	"tradeview": "https://www.bitstamp.net/market/tradeview",
}
type Link struct {
	ShortLink string
	LongLink string
	Description string
	Owner string
	New bool
}
func LoadLink(p string) Link {
	v, ok := links[p]
	link := Link{
		ShortLink: p,
		LongLink: v,
		Description: "",
		Owner: "everyone",
		New: !ok,
	}
	return link
}

type Request struct {
	Edit bool
	ShortLink string
	Url *url.URL
}

func parseRequest(path string) (Request, error) {
	u, err := url.Parse(path)
	if err != nil {
		return Request{}, fmt.Errorf("invalid request path: %s", err)
	}
	link := ""
	if strings.Contains(u.Path, "/") {
		parts := strings.SplitN(u.Path, "/", 2)
		link = parts[0]
		u.Path = parts[1]
	} else {
		link = u.Path
		u.Path = ""
	}
	r := Request{
		Edit: false,
		ShortLink: link,
		Url: u,
	}
	return r, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		createOrUpdateLink(r)
	}
	path := r.URL.Path[1:]
	edit := false
	if len(path) > 0 && path[0] == '~' {
		edit = true
		path = path[1:]
	}
	link := LoadLink(path)
	if edit {
		serveEditPage(w, link)
	} else if link.New {
		// If the link doesn't exist yet redirect to edit page
		http.Redirect(w, r, "~" + link.ShortLink, 301)
	} else {
		// This is main functionality. Redirect my dear!
		http.Redirect(w, r, link.LongLink, 301)
	}
}

func createOrUpdateLink(r *http.Request) error {
	if err := r.ParseForm(); err != null {
		return fmt.Errorf("unable to parse link form: %s", err)
	}
	// TODO: validate short link
	short := r.FormValue("short-link")
	long := r.FormValue("long-link")
	if _, ok := links[short]; ok {
		// Update
		links[short] = long
	} else {
		// Create
		links[short] = long
	}
}

func serveEditPage(w http.ResponseWriter, link Link) {
	t, err := template.ParseFiles("edit.html")
	if err != nil {
		fmt.Fprintf(w,"An error occured while parsing the template: %s", err)
	} else {
		t.Execute(w, link)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

