package main

import (
	"github.com/ashirko/imlogs/internal/db"
	"github.com/ashirko/imlogs/internal/server"
	"log"
)

func main() {
	dbInstance, err := db.Initialize()
	if err != nil {
		log.Panicln("DB initialization failed: ", err)
	}
	defer db.Close(dbInstance)

	err = server.Start()
	if err != nil {
		log.Panicln("Web Server start failed: ", err)
	}
	// TODO close the web server gracefully

}
