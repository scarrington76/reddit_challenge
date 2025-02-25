package main

import (
	"log"

	_ "github.com/lib/pq"
	"reddit_challenge/server"
)

func main() {
	app := server.NewApp()
	err := app.Serve()
	log.Println("error", err)

}
