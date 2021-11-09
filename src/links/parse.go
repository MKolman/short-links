package links

import (
	"fmt"
	"net/url"
	"strings"
)

const EditChar = "~"

// A way to store: /~first-part-of-link/and/the-rest?of=url#here
//              Edit^(--  ShortLink --) (--------- Url --------)
type Request struct {
	// Is this a request to edit this link?
	// I.e. did the request start with ~
	Edit bool
	// First part of the request path. AKA the short link.
	ShortLink string
	// The remainder of the url. This can include additional path,
	// get parameters or replacement of the search parameter.
	Url *url.URL
}

// ParsePath takes an input url path and parses it into separate parts:
// Mainly which short link it represents and what part of the path is
// "left over". In addition it checks if the path starts with '~'.
func ParsePath(path string) (Request, error) {
	u, err := url.Parse(strings.TrimPrefix(path, "/"))
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
		Edit:      strings.HasPrefix(link, EditChar),
		ShortLink: strings.TrimPrefix(link, EditChar),
		Url:       u,
	}
	return r, nil
}
