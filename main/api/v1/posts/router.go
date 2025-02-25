package posts

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route routes all the endpoints for posts
func Route(r *mux.Router) {

	r.Path("").
		Methods(http.MethodGet).
		Handler(toCtx(show()))
}
