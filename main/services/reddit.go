package services

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"reddit_challenge/model"
)

const (
	bearer    = "bearer "
	baseURL   = "https://oauth.reddit.com"
	userAgent = "reddit_challenge/1.0"
	subReddit = "r/AskReddit"
)

var accessToken string

type Reddit struct {
	Capacity     int
	CallInterval time.Duration
	PerCall      int
	consumed     int
	started      bool
	kill         chan bool
	m            sync.Mutex
}

// Token structure to store access token and refresh token
type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
	Scope        string    `json:"scope"`
	Expiry       time.Time // Add an expiry field
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

func (b *Reddit) Start() error {

	token, ok := os.LookupEnv("REDDIT_TOKEN")
	if !ok || strings.TrimSpace(token) == "" {
		log.Fatal("reddit token not found in env variables - refer to readme")
	}
	if b.started {
		return errors.New("reddit was already started")
	}

	ticker := time.NewTicker(b.CallInterval)
	b.started = true
	b.kill = make(chan bool, 1)

	go func() {
		for {
			select {
			case <-ticker.C:
				b.m.Lock()
				b.consumed -= b.PerCall
				if b.consumed < 0 {
					b.consumed = 0
				}
				b.m.Unlock()
			case <-b.kill:
				return
			}
		}
	}()

	return nil
}

func (b *Reddit) Stop() error {
	if !b.started {
		return errors.New("reddit was never started")
	}

	b.kill <- true

	return nil
}

func (b *Reddit) MakeCall(amt int) error {
	b.m.Lock()
	defer b.m.Unlock()

	if b.Capacity-b.consumed < amt {
		return errors.New("not enough capacity left to make add'l calls")
	}
	b.consumed += amt
	return nil
}

func RedditCall() error {
	var after string
	var count int

	for {
		time.Sleep(2 * time.Second)

		response, err := getPosts(after, count)
		if err != nil {
			return err
		}

		for _, post := range response.Data.Children {
			PostMap.Set(
				post.Data.ID, PostStats{
					User: post.Data.Author,
					Ups:  post.Data.Ups,
				},
			)
		}

		l := PostMap.Length()
		log.Println("the size of the map is: ", l)
		after = response.Data.After
		if count == l {
			count = 0
		} else {
			count = l
		}
	}

	return nil
}

// getPosts retrieves all the posts in a subreddit including any new posts
func getPosts(after string, count int) (response model.NewResponse, err error) {
	client := &http.Client{}
	var qry string
	if after != "" {
		qry += "&after=" + after
	}
	if count != 0 {
		qry += "&count=" + strconv.Itoa(count)
	}

	lwrSub := strings.ToLower(subReddit)
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
	if err := RedditCall(); err != nil {
		log.Println(err)
		return
	}
}
