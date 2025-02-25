package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"reddit_challenge/db"
	"reddit_challenge/server"
	"reddit_challenge/services"
)

func main() {
	d, err := sql.Open("postgres", dataSource())
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()
	// CORS is enabled only in prod profile
	cors := os.Getenv("profile") == "prod"
	app := server.NewApp(db.NewDB(d), cors)
	err = app.Serve()
	if err = services.RedditCall(); err != nil {
		log.Println(err)
	}
	log.Println("error", err)

}

func dataSource() string {
	host := "localhost"
	pass := "pass"
	if os.Getenv("profile") == "prod" {
		host = "db"
		pass = os.Getenv("db_pass")
	}
	return "postgresql://" + host + ":5432/goxygen" +
		"?user=goxygen&sslmode=disable&password=" + pass
}
