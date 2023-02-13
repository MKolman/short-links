package server

import (
	"fmt"
	"html/template"
	"net/http"
)

type EditViewModel struct {
	ShortLink string
	LongLink  string
	New       bool
	Error     error
}

func (s *Server) serveEditPage(w http.ResponseWriter, shortLink string, err error) {
	data := EditViewModel{
		ShortLink: shortLink,
		LongLink:  "",
		New:       true,
		Error:     err,
	}
	if link, err := s.store.Get(shortLink); err == nil {
		data.LongLink = link.LongLink
		data.New = false
	}
	s.renderEditPage(w, &data)
}

func (s *Server) renderEditPage(w http.ResponseWriter, data *EditViewModel) {
	t, err := template.ParseFiles("templates/edit.html")
	if err != nil {
		fmt.Fprintf(w, "An error occured while parsing the template: %s", err)
	} else {
		t.Execute(w, data)
	}
}
