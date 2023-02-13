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
		s.postHandler(w, r)
		return
	}
	var renderErr error
	parsed, err := links.ParsePath(fmt.Sprintf("%s?%s#%s", r.URL.Path, r.URL.RawQuery, r.URL.Fragment))
	if err != nil {
		renderErr := fmt.Errorf("unable to parse path: %w", err)
		log.Errorf(renderErr.Error())
	}
	if parsed.Edit {
		s.serveEditPage(w, parsed.ShortLink, renderErr)
	} else if link, err := s.store.Get(parsed.ShortLink); err == nil {
		// This is main functionality. Redirect my dear!
		url, err := links.Merge(link.LongLink, parsed.Url)
		if err != nil {
			url = link.LongLink
		}
		log.Infof("Redirecting %q -> %s", parsed.ShortLink, url)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		// If the link doesn't exist yet redirect to edit page
		log.Warnf("Link %q does not exist. Redirecting to edit page.", parsed.ShortLink)
		http.Redirect(w, r, "/~"+parsed.ShortLink, http.StatusFound)
	}
}

func (s *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	data := EditViewModel{
		ShortLink: "",
		LongLink:  "",
		New:       true,
		Error:     nil,
	}
	if err := r.ParseForm(); err != nil {
		data.Error = fmt.Errorf("unable to parse form: %w", err)
		s.renderEditPage(w, &data)
		return
	}
	data.ShortLink = r.FormValue("short-link")
	data.LongLink = r.FormValue("long-link")
	if err := validation.ValidateShort(data.ShortLink); err != nil {
		data.Error = fmt.Errorf("invalid short link: %w", err)
		s.renderEditPage(w, &data)
		return
	}
	if err := validation.ValidateLong(data.LongLink); err != nil {
		data.Error = fmt.Errorf("invalid long link: %w", err)
		s.renderEditPage(w, &data)
		return
	}

	if l, err := s.store.Get(data.LongLink); err == nil {
		// Update
		log.Infof("Update link %q", l.ShortLink)
		l.LongLink = data.LongLink
		if err := s.store.Update(l); err != nil {
			data.Error = err
			s.renderEditPage(w, &data)
			return
		}
		data.New = false
	} else {
		// Create
		l := &db.Link{
			ShortLink: data.ShortLink,
			LongLink:  data.LongLink,
		}
		log.Infof("Create link %q -> %s", l.ShortLink, l.LongLink)
		_, err := s.store.Create(l)
		if err != nil {
			data.Error = err
			s.renderEditPage(w, &data)
			return
		}
		data.New = false
	}
	s.renderEditPage(w, &data)
}

func Run(store db.Db, port int) {
	server := Server{store: store}
	http.HandleFunc("/", server.handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
