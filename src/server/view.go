package server

import (
	"fmt"
	"html/template"
	"net/http"
)

type EditViewModel struct {
	ShortLink   string
	LongLink    string
	Description string
	New         bool
}

func (s *Server) serveEditPage(w http.ResponseWriter, shortLink string) {
	data := EditViewModel{
		ShortLink:   shortLink,
		LongLink:    "",
		Description: "",
		New:         true,
	}
	if link, err := s.store.Get(shortLink); err == nil {
		data.LongLink = link.LongLink
		data.Description = link.Description
		data.New = false
	}

	t, err := template.ParseFiles("templates/edit.html")
	if err != nil {
		fmt.Fprintf(w, "An error occured while parsing the template: %s", err)
	} else {
		t.Execute(w, data)
	}
}
