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
		validateTestCase{"t", true},
		validateTestCase{"basic", true},
		validateTestCase{"long-link_withmany.chars", true},
		validateTestCase{"", false},
		validateTestCase{"~basic", false},
		validateTestCase{"a/t", false},
		validateTestCase{"m:a", false},
	}
	for _, c := range cases {
		err := ValidateShort(c.input)
		assert.Equal(t, c.valid, err == nil)
	}
}

func TestValidateLong(t *testing.T) {
	cases := []validateTestCase{
		validateTestCase{"https://www.kolman.si", true},
		validateTestCase{"ssh://a", true},
		validateTestCase{"http://sub.domain.com/help?a=23#15", true},
		validateTestCase{"", false},
		validateTestCase{"test", false},
		validateTestCase{"kolman.si", false},
		validateTestCase{"_://_", false},
	}
	for _, c := range cases {
		err := ValidateLong(c.input)
		assert.Equal(t, c.valid, err == nil, fmt.Errorf("%s: %t %s", c.input, c.valid, err))
	}
}
