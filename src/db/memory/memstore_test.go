package memory

import (
	"net/url"
	"short-links/db/internal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemDb(t *testing.T) {
	mdb, err := CreateMemDb(&url.URL{})
	assert.Nil(t, err)
	internal.TestDbAll(t, mdb)
}
