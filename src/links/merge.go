package links

import (
	"fmt"
	"net/url"
	"path"
)

// Append the url parameters from extraUrl onto longLink
func Merge(longLink string, extraUrl *url.URL) (string, error) {
	r, err := url.Parse(longLink)
	if err != nil {
		return "", fmt.Errorf("error parsing longLink: %s", err)
	}
	r.Path = path.Join(r.Path, extraUrl.Path)
	q := r.Query()
	for k, vs := range extraUrl.Query() {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	r.RawQuery = q.Encode()
	if len(extraUrl.Fragment) > 0 {
		r.Fragment = extraUrl.Fragment
	}
	return r.String(), nil
}
