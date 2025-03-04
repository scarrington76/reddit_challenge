package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/time/rate"
	"reddit_challenge/model"
)

const (
	bearer    = "bearer "
	baseURL   = "https://oauth.reddit.com"
	userAgent = "reddit_challenge/1.0"

	// redditUsed is the header to identify the number of requests used in the period
	redditUsed = "X-Ratelimit-Used"

	// redditRemaining is the header to identify the number of requests remaining in the period
	redditRemaining = "X-Ratelimit-Remaining"

	// redditRemaining is the header to identify the time (in seconds) until the period reset
	redditReset = "X-Ratelimit-Reset"

	// 100 queries per minute (QPM) per OAuth client id
	// QPM limits will be an average over a time window (currently 10 minutes) to support bursting requests.
	redditEventsPerSecond = 100 / 60
)

type RedditRequest struct {
	Sub   string
	Map   *SafeMap
	after string
	count int
}

// RedditErrorResponse is an error response structure
type RedditErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

var PostMap *SafeMap

var Subs []RedditRequest

func Start(tok string) error {

	// set rate limiter to the number of events that reddit allows. Burst is set to
	// five; could not find anywhere in reddit documents where max burst is set
	limiter := rate.NewLimiter(rate.Limit(redditEventsPerSecond/len(Subs)), 5)

	for _, req := range Subs {
		go func() {
			err := RedditCall(req, tok, limiter)
			if err != nil {
				log.Printf("error processing %s: %v", req.Sub, err)
				return
			}
		}()
	}

	return nil
}

func RedditCall(sub RedditRequest, tok string, limit *rate.Limiter) error {

	for {
		ctx := context.Background()
		err := limit.Wait(ctx)
		if err != nil {
			fmt.Println("Error waiting for token:", err)
			return err
		}
		response, err := getPosts(sub, tok)
		if err != nil {
			log.Print("error getting posts: ", err)
			return err
		}

		for _, post := range response.Data.Children {
			sub.Map.Set(
				post.Data.ID, PostStats{
					User:  post.Data.Author,
					Ups:   post.Data.Ups,
					Title: post.Data.Title,
				},
			)
		}

		l := sub.Map.Length()
		log.Println("the size of the "+sub.Sub+" map is: ", l)
		sub.after = response.Data.After
		if sub.count == l {
			sub.count = 0
		} else {
			sub.count = l
		}
	}
}

// getPosts retrieves all the posts in a subreddit including any new posts and places them into a map
// for storage purposes
func getPosts(sub RedditRequest, tok string) (response model.NewResponse, err error) {
	client := &http.Client{}
	var qry string

	if sub.after != "" {
		qry += "&after=" + sub.after
	}

	if sub.count != 0 {
		qry += "&count=" + strconv.Itoa(sub.count)
	}

	lwrSub := strings.ToLower(sub.Sub)
	req, err := http.NewRequest(
		"GET", baseURL+"/"+lwrSub+"/new?limit=100"+qry, nil,
	)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", bearer+tok)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	used := resp.Header.Get(redditUsed)
	remaining := resp.Header.Get(redditRemaining)
	reset := resp.Header.Get(redditReset)
	log.Printf("used: %v, remaining: %v, reset: %v\n", used, remaining, reset)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &response); err != nil {
		var e *json.SyntaxError
		if errors.As(err, &e) {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		return
	}

	return
}

// TrackSub creates an object to hold information for the reddit sub the client would like to track
func TrackSub(sub string) {
	req := RedditRequest{
		Sub: sub,
		Map: NewSafeMap(),
	}
	Subs = append(Subs, req)
}
