package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type validateTestCase struct {
	input string
	valid bool
}

func TestValidateShort(t *testing.T) {
	cases := []validateTestCase{
		{"t", true},
		{"basic", true},
		{"long-link_withmany.chars", true},
		{"", false},
		{"~basic", false},
		{"a/t", false},
		{"m:a", false},
	}
	for _, c := range cases {
		err := ValidateShort(c.input)
		assert.Equal(t, c.valid, err == nil)
	}
}

func TestValidateLong(t *testing.T) {
	cases := []validateTestCase{
		{"https://www.kolman.si", true},
		{"ssh://a", true},
		{"http://sub.domain.com/help?a=23#15", true},
		{"", false},
		{"test", false},
		{"kolman.si", false},
		{"_://_", false},
	}
	for _, c := range cases {
		err := ValidateLong(c.input)
		assert.Equal(t, c.valid, err == nil, fmt.Errorf("%s: %t %s", c.input, c.valid, err))
	}
}
