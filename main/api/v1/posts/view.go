package posts

import (
	"net/http"

	"reddit_challenge/helpers"
)

// show displays the information for posts (post name and # of upvotes) in descending order
func show() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			posts := ctx.Value("posts").([]response)

			dto := []any{}
			for _, p := range posts {
				l := struct {
					Name  string `json:"name"`
					Votes int    `json:"votes"`
				}{
					Name:  p.Name,
					Votes: p.Votes,
				}
				dto = append(dto, l)
			}

			helpers.WriteJSON(w, dto)
		},
	)
}
