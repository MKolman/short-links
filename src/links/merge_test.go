package links

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MergeTestCase struct {
	link string
	url  url.URL
	want string
}

func TestMerge(t *testing.T) {
	cases := []MergeTestCase{
		{
			link: "http://kolman.si",
			url:  url.URL{},
			want: "http://kolman.si",
		},
		{
			link: "https://www.kolman.si/test?j=3#hash",
			url:  url.URL{},
			want: "https://www.kolman.si/test?j=3#hash",
		},
		{
			link: "http://kolman.si",
			url:  url.URL{Path: "test", RawQuery: "what=1"},
			want: "http://kolman.si/test?what=1",
		},
		{
			link: "http://kolman.si/asdf?what=0#removeme",
			url:  url.URL{Path: "test", RawQuery: "what=1", Fragment: "newhash"},
			want: "http://kolman.si/asdf/test?what=0&what=1#newhash",
		},
	}
	for _, tc := range cases {
		got, err := Merge(tc.link, &tc.url)
		assert.Nil(t, err)
		assert.Equal(t, tc.want, got)
	}
}
