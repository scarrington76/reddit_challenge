package api

import (
	"github.com/gorilla/mux"
	v1 "reddit_challenge/api/v1"
)

// Route routes all the endpoints for the api
func Route(r *mux.Router) {

	l := r.PathPrefix("/v1").Subrouter()
	v1.Route(l)
}
