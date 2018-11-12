package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("Kontaktmanagement")
	db, err := New()
	if err != nil {
		log.Fatalf("%v", err)
	}
	InitAndStartRouter(db)
}
