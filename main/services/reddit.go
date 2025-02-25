package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/time/rate"
	"reddit_challenge/model"
)

const (
	bearer    = "bearer "
	baseURL   = "https://oauth.reddit.com"
	userAgent = "reddit_challenge/1.0"
	subReddit = "r/AskReddit"

	// redditUsed is the header to identify the number of requests used in the period
	redditUsed = "X-Ratelimit-Used"

	// redditRemaining is the header to identify the number of requests remaining in the period
	redditRemaining = "X-Ratelimit-Remaining"

	// redditRemaining is the header to identify the time (in seconds) until the period reset
	redditReset = "X-Ratelimit-Reset"

	// 100 queries per minute (QPM) per OAuth client id
	// QPM limits will be an average over a time window (currently 10 minutes) to support bursting requests.
)

var accessToken string

type SubReddit struct {
	Name  string
	after string
	count int
}

// RedditErrorResponse is an error response structure
type RedditErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

var PostMap *SafeMap

func init() {
	PostMap = NewSafeMap()
}

func Start() error {
	if accessToken == "" || strings.TrimSpace(accessToken) == "" {
		log.Println("reddit service cannot being; access token not present; will re-attempt in 20 seconds")
		time.Sleep(time.Second * 20)
		Start()
	}
	// TODO: get subreddits from config
	sub := SubReddit{
		Name: subReddit,
	}
	go func() {
		err := RedditCall(sub)
		if err != nil {
			log.Printf("error processing %s: %v", sub.Name, err)
			return
		}
	}()

	return nil
}

func RedditCall(sub SubReddit) error {

	limiter := rate.NewLimiter(rate.Limit(100/60), 5)

	for {
		ctx := context.Background()
		err := limiter.Wait(ctx)
		if err != nil {
			fmt.Println("Error waiting for token:", err)
			return err
		}
		response, err := getPosts(sub)
		if err != nil {
			return err
		}

		for _, post := range response.Data.Children {
			PostMap.Set(
				post.Data.ID, PostStats{
					User:  post.Data.Author,
					Ups:   post.Data.Ups,
					Title: post.Data.Title,
				},
			)
		}

		l := PostMap.Length()
		log.Println("the size of the map is: ", l)
		sub.after = response.Data.After
		if sub.count == l {
			sub.count = 0
		} else {
			sub.count = l
		}
	}

	return nil
}

// getPosts retrieves all the posts in a subreddit including any new posts and places them into a map
// for storage purposes
func getPosts(sub SubReddit) (response model.NewResponse, err error) {
	client := &http.Client{}
	var qry string

	if sub.after != "" {
		qry += "&after=" + sub.after
	}

	if sub.count != 0 {
		qry += "&count=" + strconv.Itoa(sub.count)
	}

	lwrSub := strings.ToLower(sub.Name)
	req, err := http.NewRequest(
		"GET", baseURL+"/"+lwrSub+"/new?limit=100"+qry, nil,
	)
	if err != nil {
		return
	}

	req.Header.Set("Authorization", bearer+accessToken)
	req.Header.Set("User-Agent", userAgent)

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

	return
}

func SetAccessToken(tok string) {
	if tok == "" || strings.TrimSpace(tok) == "" {
		log.Println("reddit access token not set; token is empty")
		return
	}
	accessToken = tok
	log.Println("access token set to: ", accessToken)
}
