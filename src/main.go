package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"short-links/db"
	"short-links/links"
	"short-links/validation"

	_ "short-links/db/memory"
)

var (
	dbConnection = flag.String("db-connection", "mem://local", "Connection URI string.")
	port         = flag.Int("port", 8081, "Port on which to serve.")
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		createOrUpdateLink(r)
	}
	parsed, _ := links.ParsePath(fmt.Sprintf("%s?%s#%s", r.URL.Path, r.URL.RawQuery, r.URL.Fragment))
	if parsed.Edit {
		ServeEditPage(w, parsed.ShortLink)
	} else if link, err := db.Store.Get(parsed.ShortLink); err == nil {
		// This is main functionality. Redirect my dear!
		url, _ := links.Merge(link.LongLink, parsed.Url)
		http.Redirect(w, r, url, 301)
	} else {
		// If the link doesn't exist yet redirect to edit page
		http.Redirect(w, r, "/~"+parsed.ShortLink, 301)
	}
}

func createOrUpdateLink(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("unable to parse form: %s", err)
	}
	short := r.FormValue("short-link")
	long := r.FormValue("long-link")
	description := r.FormValue("description")
	if err := validation.ValidateShort(short); err != nil {
		return fmt.Errorf("invalid short link: %s", err)
	}
	if err := validation.ValidateLong(long); err != nil {
		return fmt.Errorf("invalid long link: %s", err)
	}

	if l, err := db.Store.Get(short); err == nil {
		// Update
		l.LongLink = long
		if err := db.Store.Update(l); err != nil {
			return err
		}
	} else {
		// Create
		l := &db.Link{
			ShortLink:   short,
			LongLink:    long,
			Description: description,
		}
		_, err := db.Store.Create(l)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	err := db.LoadDb(*dbConnection)
	db.Store.Create(&db.Link{ShortLink: "kolman", LongLink: "https://www.kolman.si"})
	if err != nil {
		panic(fmt.Sprintf("unable connect to database: %s", err))
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
