package posts_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"reddit_challenge/server/test"
)

// getPostsFail tests the bad path(s) for getting posts
func getPostsFail(t *testing.T) {
	results := test.NewRequestNoBody(t, http.MethodGet, "/api/v1/posts", http.StatusOK)

	if len(results) > 10 {
		t.Fatal("length of post results is more than 10")
	}
}

// getPosts tests the happy path for getting posts
func getPostsSuccessful(t *testing.T) {
	results := test.NewRequestNoBody(t, http.MethodGet, "/api/v1/posts", http.StatusOK)

	if len(results) == 0 {
		t.Fatal("posts are empty")
	}

	var mx float64
	for index, result := range results {
		name, ok := result["name"].(string)
		if !ok || strings.TrimSpace(name) == "" {
			t.Fatal("post name isn't a string, or is an empty string val")
		}
		votes, ok := result["votes"].(float64)
		if !ok || votes < 0 {
			t.Fatal("number of posts isn't a float64, or is a negative val")
		}
		if index == 0 {
			mx = votes
			continue // set initial number of votes
		}
		if votes > mx {
			t.Fatal("returned json is not sorted by # of votes")
		} else {
			mx = votes
		}
	}
}

// TestPostEndpoints tests each of the posts endpoints
func TestPostEndpoints(t *testing.T) {

	test.Setup(t)

	// immediately run the reddit tracking service; we'll check for more than 10 posts which
	// should fail given that this is running immediately (and doesn't have time to collect
	// more than 10 posts from our tracker)
	getPostsFail(t)

	// sleep the system and let the sub collect posts
	time.Sleep(10 * time.Second)
	getPostsSuccessful(t)

}
