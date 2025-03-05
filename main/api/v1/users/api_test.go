package users_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"reddit_challenge/server/test"
)

// getUsersFail tests the bad path(s) for getting users
func getUsersFail(t *testing.T) {
	results := test.NewRequestNoBody(t, http.MethodGet, "/api/v1/users", http.StatusOK)

	if len(results) > 10 {
		t.Fatal("length of user results is more than 10")
	}
}

// getUsersSuccessful tests the happy path for getting users
func getUsersSuccessful(t *testing.T) {
	results := test.NewRequestNoBody(t, http.MethodGet, "/api/v1/users", http.StatusOK)

	if len(results) == 0 {
		t.Fatal("posts are empty")
	}

	var mx float64
	for index, result := range results {
		user, ok := result["username"].(string)
		if !ok || strings.TrimSpace(user) == "" {
			t.Fatal("username isn't a string, or is an empty string val")
		}
		posts, ok := result["posts"].(float64)
		if !ok || posts <= 0 {
			t.Fatal("number of posts isn't a float64, or is a zero / negative val")
		}
		if index == 0 {
			mx = posts
			continue // set initial number of posts
		}
		if posts > mx {
			t.Fatal("returned json is not sorted by # of posts")
		} else {
			mx = posts
		}
	}
}

// TestUserEndpoints tests each of the users endpoints
func TestUserEndpoints(t *testing.T) {

	test.Setup(t)

	// immediately run the reddit tracking service; we'll check for more than 10 users which
	// should fail given that this is running immediately (and doesn't have time to collect
	// more than 10 posts from our tracker)
	getUsersFail(t)

	// sleep the system and let the sub collect users
	time.Sleep(10 * time.Second)
	getUsersSuccessful(t)

}
