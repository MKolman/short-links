package validation

import (
	"fmt"
	"net/url"
)

var BlockedChars string = "~:/?#[]@!$&'()*+,;=\"<>%{}|\\^`"

func ValidateShort(link string) error {
	if len(link) == 0 {
		return fmt.Errorf("link has to be at least one character long")
	}
	for i, c := range link {
		for _, b := range BlockedChars {
			if c == b {
				return fmt.Errorf("short link contains character %c at position %d", c, i)
			}
		}
	}
	return nil
}

func ValidateLong(link string) error {
	u, err := url.Parse(link)
	if err != nil {
		return fmt.Errorf("url %q cannot be parsed: %w", link, err)
	}
	if u.Scheme == "" {
		return fmt.Errorf("url scheme not defined in %q", link)
	}
	if u.Host == "" {
		return fmt.Errorf("url domain not defined in %q", link)
	}

	return nil
}
