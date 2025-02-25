package v1

import (
	"github.com/gorilla/mux"
	"reddit_challenge/api/v1/posts"
	"reddit_challenge/api/v1/system"
	"reddit_challenge/api/v1/users"
)

// Route routes all the endpoints for v1 of the api
func Route(r *mux.Router) {

	l := r.PathPrefix("/posts").Subrouter()
	posts.Route(l)

	u := r.PathPrefix("/users").Subrouter()
	users.Route(u)

	s := r.PathPrefix("/system").Subrouter()
	system.Route(s)
}
