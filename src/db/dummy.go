package db

// Dummy type with Db interface that panics
// on all operations.
type dummyDb struct{}

func (db *dummyDb) Get(_ string) (*Link, error) {
	panic("The database has not been loaded yet. Plase call LoadDb before making any calls.")
}

func (db *dummyDb) Create(_ *Link) (*Link, error) {
	panic("The database has not been loaded yet. Plase call LoadDb before making any calls.")
}

func (db *dummyDb) Update(_ *Link) error {
	panic("The database has not been loaded yet. Plase call LoadDb before making any calls.")
}
