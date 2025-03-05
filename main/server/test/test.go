package test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"reddit_challenge/server"
	"reddit_challenge/services"
)

// testApp is the instance of the app used by the test cases
var testApp server.App

// Setup creates the test server with a reddit access token
func Setup(t *testing.T) {

	// token must be set manually for testing
	tok := ""
	t.Setenv("ACCESS_TOKEN", tok)
	redditToken, ok := os.LookupEnv("ACCESS_TOKEN")
	if !ok || strings.TrimSpace(redditToken) == "" {
		t.Fatal("reddit access token required - set token variable in /main/server/test/test.go and retry test")
	}
	services.TrackSub("r/askreddit")
	err := services.Start(redditToken)
	if err != nil {
		log.Panic("services error: ", err)
	}
	testApp = server.NewApp(redditToken)
}

// NewRequestNoBody creates and executes an http request on the test server and returns the response
func NewRequestNoBody(t *testing.T, method string, url string, expectedCode int) []map[string]interface{} {
	req, _ := http.NewRequest(method, url, nil)
	rr := httptest.NewRecorder()
	testApp.Router.ServeHTTP(rr, req)
	if rr.Code != expectedCode {
		t.Fatalf("expected response code %d but received response code %d", expectedCode, rr.Code)
	}

	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}

	var m []map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		t.Fatal(err)
	}

	return m

}
