package main

import (
	"flag"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"reddit_challenge/server"
	"reddit_challenge/services"
)

func main() {

	redditToken, ok := os.LookupEnv("ACCESS_TOKEN")
	if ok && strings.TrimSpace(redditToken) != "" {
		log.Println("reddit access token found")
	} else {
		log.Fatal("reddit token required - set env variable for ACCESS_TOKEN and re-start app")
	}

	redditSub := flag.String("sub", "r/askreddit", "sub-reddit to subscribe to")
	if redditSub == nil {
		log.Panic("no reddit sub provided in command line arguments (i.e. go build main.go -sub r/askreddit)")
	}

	services.TrackSub(*redditSub)
	err := services.Start(redditToken)
	if err != nil {
		log.Panic("services error: ", err)
	}
	app := server.NewApp(redditToken)
	err = app.Serve()
	log.Println("error", err)

}
