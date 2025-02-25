package system

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route routes the endpoints for the system endpoints
func Route(r *mux.Router) {

	r.Path("/callback").
		Methods(http.MethodGet, http.MethodPut, http.MethodPost).
		Handler(handleCallback())
}
