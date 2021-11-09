package memory

import (
	"fmt"
	"net/url"
	"short-links/db"
)

type MemDb struct {
	store map[string]db.Link
}

func (m *MemDb) Get(shortLink string) (*db.Link, error) {
	l, ok := m.store[shortLink]
	if !ok {
		return nil, fmt.Errorf("short link does not exist")
	}
	return &l, nil
}

func (m *MemDb) Create(l *db.Link) (*db.Link, error) {
	l.Id = l.ShortLink
	m.store[l.ShortLink] = *l
	return l, nil
}

func (m *MemDb) Update(l *db.Link) error {
	if l.Id != l.ShortLink {
		return fmt.Errorf("cannot change short link of an existing object")
	}
	_, err := m.Create(l)
	return err
}

func CreateMemDb(_ *url.URL) (db.Db, error) {
	return &MemDb{
		store: make(map[string]db.Link),
	}, nil
}

func init() {
	db.RegisterDbBuilder("mem", CreateMemDb)
	db.RegisterDbBuilder("memory", CreateMemDb)
}
