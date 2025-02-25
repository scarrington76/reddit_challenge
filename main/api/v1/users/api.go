package users

import (
	"context"
	"net/http"
	"sort"

	"reddit_challenge/services"
)

// toCtx places a slice of the users and their # of posts into context
func toCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			usermap := make(map[string]int)
			users := []user{}

			posts := services.ClonePostMap()
			for _, stats := range posts {
				usr, ok := usermap[stats.User]
				if ok {
					usermap[stats.User] = usr + 1
				} else {
					usermap[stats.User] = 1
				}
			}

			for k, v := range usermap {
				users = append(
					users, user{
						username: k,
						posts:    v,
					},
				)
			}

			sort.Slice(
				users, func(i, j int) bool {
					return users[i].posts > users[j].posts
				},
			)

			ctx = context.WithValue(ctx, "users", users)

			h.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
