package db

import (
	"fmt"
	"net/url"
)

var Store Db = &dummyDb{}

var builders map[string]func(*url.URL) (Db, error)

type Link struct {
	Id          interface{}
	ShortLink   string
	LongLink    string
	Description string
}

type Db interface {
	Get(string) (*Link, error)
	Create(*Link) (*Link, error)
	Update(*Link) error
}

func RegisterDbBuilder(scheme string, builder func(*url.URL) (Db, error)) {
	if builders == nil {
		builders = make(map[string]func(*url.URL) (Db, error))
	}
	builders[scheme] = builder
}

func LoadDb(connectionUri string) error {
	con, err := url.Parse(connectionUri)
	if err != nil {
		return fmt.Errorf("unable to parse connection string %q: %s", connectionUri, err)
	}
	if builder, ok := builders[con.Scheme]; ok {
		var err error
		Store, err = builder(con)
		return err
	}
	return fmt.Errorf("no database handler implemented for scheme %q", con.Scheme)
}
