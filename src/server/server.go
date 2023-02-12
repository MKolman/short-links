package server

import (
	"fmt"
	"net/http"
	"short-links/db"
	"short-links/links"
	"short-links/validation"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	store db.Db
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := s.createOrUpdateLink(r); err != nil {
			log.Error(err)
		}
	}
	parsed, _ := links.ParsePath(fmt.Sprintf("%s?%s#%s", r.URL.Path, r.URL.RawQuery, r.URL.Fragment))
	if parsed.Edit {
		s.serveEditPage(w, parsed.ShortLink)
	} else if link, err := s.store.Get(parsed.ShortLink); err == nil {
		// This is main functionality. Redirect my dear!
		url, _ := links.Merge(link.LongLink, parsed.Url)
		log.Infof("Redirecting %q -> %s", parsed.ShortLink, url)
		http.Redirect(w, r, url, 301)
	} else {
		// If the link doesn't exist yet redirect to edit page
		log.Warnf("Link %q does not exist. Redirecting to edit page.", parsed.ShortLink)
		http.Redirect(w, r, "/~"+parsed.ShortLink, 301)
	}
}

func (s *Server) createOrUpdateLink(r *http.Request) error {
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

	if l, err := s.store.Get(short); err == nil {
		// Update
		log.Infof("Update link %q", l.ShortLink)
		l.LongLink = long
		if err := s.store.Update(l); err != nil {
			return err
		}
	} else {
		// Create
		l := &db.Link{
			ShortLink:   short,
			LongLink:    long,
			Description: description,
		}
		log.Infof("Create link %q -> %s", l.ShortLink, l.LongLink)
		_, err := s.store.Create(l)
		if err != nil {
			return err
		}
	}
	return nil
}

func Run(store db.Db, port int) {
	server := Server{store: store}
	http.HandleFunc("/", server.handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
