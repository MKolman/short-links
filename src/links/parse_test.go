package links

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ParseTestCase struct {
	inp  string
	want *Request
}

func TestParsePath(t *testing.T) {
	data := []ParseTestCase{
		{
			inp: "/testme/",
			want: &Request{
				Edit:      false,
				ShortLink: "testme",
				Url:       &url.URL{},
			},
		},
		{
			inp: "/~testme/",
			want: &Request{
				Edit:      true,
				ShortLink: "testme",
				Url:       &url.URL{},
			},
		},
		{
			inp: "/testme/test2",
			want: &Request{
				Edit:      false,
				ShortLink: "testme",
				Url:       &url.URL{Path: "test2"},
			},
		},
		{
			inp: "/~test-me_/test2?q=t&b=x#hashme",
			want: &Request{
				Edit:      true,
				ShortLink: "test-me_",
				Url: &url.URL{
					Path:     "test2",
					Fragment: "hashme",
					RawQuery: "q=t&b=x",
				},
			},
		},
	}
	for _, tc := range data {
		got, err := ParsePath(tc.inp)
		assert.Nil(t, err)
		assert.Equal(t, *tc.want, got)
	}

}
