package internal

import (
	"short-links/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDbAll(t *testing.T, store db.Db) {
	TestDbGetEmpty(t, store)
	TestDbCreateAndGet(t, store)
	TestDbCreateAndWrongUpdate(t, store)
	TestDbCreateUpdateGet(t, store)
}

func TestDbGetEmpty(t *testing.T, store db.Db) {
	_, err := store.Get("this-key-does-not-exist")
	assert.NotNil(t, err)
}

func TestDbCreateAndGet(t *testing.T, store db.Db) {
	l, err := store.Create(&db.Link{ShortLink: "test", LongLink: "http://a"})
	assert.Nil(t, err)
	assert.NotNil(t, l.Id)
	assert.Equal(t, "test", l.ShortLink)
	assert.Equal(t, "http://a", l.LongLink)

	l2, err := store.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, l, l2)
}

func TestDbCreateAndWrongUpdate(t *testing.T, store db.Db) {
	l, err := store.Create(&db.Link{ShortLink: "test", LongLink: "http://a"})
	assert.Nil(t, err)

	l.ShortLink = "x"
	err = store.Update(l)
	assert.NotNil(t, err)
}

func TestDbCreateUpdateGet(t *testing.T, store db.Db) {
	l, err := store.Create(&db.Link{ShortLink: "test", LongLink: "http://a"})
	assert.Nil(t, err)

	l.LongLink = "https://b"
	err = store.Update(l)
	assert.Nil(t, err)

	l2, err := store.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, l2, l)
}
