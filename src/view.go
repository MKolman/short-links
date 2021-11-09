package main

import (
	"fmt"
	"html/template"
	"net/http"
	"short-links/db"
)

type EditViewModel struct {
	ShortLink   string
	LongLink    string
	Description string
	New         bool
}

func ServeEditPage(w http.ResponseWriter, shortLink string) {
	data := EditViewModel{
		ShortLink:   shortLink,
		LongLink:    "",
		Description: "",
		New:         true,
	}
	if link, err := db.Store.Get(shortLink); err == nil {
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
