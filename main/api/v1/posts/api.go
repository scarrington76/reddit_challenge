package posts

import (
	"context"
	"net/http"
	"sort"

	"reddit_challenge/model"
	"reddit_challenge/services"
)

// toCtx obtains the posts, sorts them by desc and places them into ctx
func toCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			posts := []model.Post{}
			pmap := services.ClonePostMap()

			for _, v := range pmap {
				posts = append(
					posts, model.Post{
						Name:  v.Title,
						Votes: v.Ups,
					},
				)
			}

			sort.Slice(
				posts, func(i, j int) bool {
					return posts[i].Votes > posts[j].Votes
				},
			)

			ctx = context.WithValue(ctx, "posts", posts)
			h.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
