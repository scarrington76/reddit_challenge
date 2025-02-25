package users

import (
	"net/http"

	"reddit_challenge/helpers"
)

// show displays the users and votes to the client
func show() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			users := ctx.Value("users").([]user)

			dto := []any{}
			for _, u := range users {
				l := struct {
					User  string `json:"username"`
					Posts int    `json:"posts"`
				}{
					User:  u.username,
					Posts: u.posts,
				}
				dto = append(dto, l)
			}

			helpers.WriteJSON(w, dto)
		},
	)
}
