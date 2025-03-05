package services

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/time/rate"
	"reddit_challenge/model"
)

const (
	// bearer is used to form the bearer token auth
	bearer = "bearer "

	// baseURL is the base url for all calls to the reddit api
	baseURL = "https://oauth.reddit.com"

	// userAgent is provided to reddit to identify where the request is coming from
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

// RedditRequest holds the information associated to a reddit sub that is tracked by the system
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

// Subs holds all the current reddit requests. For initial purposes, only 1 RedditRequest
// is currently supported. Future revisions would allow for multiple threads.
var Subs []RedditRequest

// Start kicks off the service to track a reddit thread based on reddit's API
// limits. At this time only 1 subreddit thread is supported. Future iterations
// would include multiple threads and a restart() whenever changes occurred in the
// number of threads tracked
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

// RedditCall maintains the polling to the reddit API and posts the returned information to the respective map
func RedditCall(sub RedditRequest, tok string, limit *rate.Limiter) error {

	for {
		ctx := context.Background()
		err := limit.Wait(ctx)
		if err != nil {
			log.Println("error waiting for token: ", err)
			return err
		}
		response, err := getPosts(sub, tok)
		if err != nil {
			log.Println("error getting posts: ", err)
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

// getPosts executes a call to the reddit API returns the map for processing
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
	req, err := http.NewRequest(http.MethodGet, baseURL+"/"+lwrSub+"/new?limit=100"+qry, nil)
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
		return
	}

	if len(response.Data.Children) == 0 {
		log.Println("warning - response is empty")
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
