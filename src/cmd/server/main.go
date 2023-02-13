package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"short-links/db"
	"short-links/server"

	_ "short-links/db/memory"
	_ "short-links/db/redis"
)

var (
	dbConnection = flag.String("db-connection", "mem://local", "Connection URI string.")
	port         = flag.Int("port", 8081, "Port on which to serve.")
)

func main() {
	flag.Parse()
	log.Info("Application starting...")
	log.Infof("Loading store from %q.", *dbConnection)
	store, err := db.LoadDb(*dbConnection)
	if err != nil {
		log.Fatalf("Unable connect to database: %s", err)
	}
	store.Create(&db.Link{ShortLink: "kolman", LongLink: "https://www.kolman.si"})
	server.Run(store, *port)
}
