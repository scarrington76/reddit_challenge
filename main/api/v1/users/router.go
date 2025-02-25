package users

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route routes all the endpoints for the users endpoints
func Route(r *mux.Router) {

	r.Path("").
		Methods(http.MethodGet).
		Handler(toCtx(show()))
}
