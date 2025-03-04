package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"reddit_challenge/api"
)

type App struct {
	Router      http.Handler
	RedditToken string
}

// NewApp creates an instance of the reddit application
func NewApp(token string, sub string) App {

	// create the router & subroutes
	r := mux.NewRouter()
	a := r.PathPrefix("/api").Subrouter()
	api.Route(a)

	// apply CORS to handle running the app locally. Would skip in prod.
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)(r)

	app := App{
		Router:      corsHandler,
		RedditToken: token,
	}
	return app
}

// Serve builds the server object and instantiates the web server
func (a *App) Serve() error {
	port := ":8080"
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("api is available on port %s\n", port)
	return srv.ListenAndServe()
}
