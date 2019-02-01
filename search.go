package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
   	"regexp"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/google/uuid"
)

const (
	magnitudeThreshold     = 0.1
)

var (
	// twitter API secrets
	consumerKey = os.Getenv("T_CONSUMER_KEY")
	consumerSecret = os.Getenv("T_CONSUMER_SECRET")
	accessToken = os.Getenv("T_ACCESS_TOKEN")
	accessSecret = os.Getenv("T_ACCESS_SECRET")

	// validation expressions
	userReg = regexp.MustCompile(`@[\w]*`)
	nonCharReg = regexp.MustCompile(`[^a-zA-Z#]`)
	shortReg = regexp.MustCompile(`\b[a-z]{1,2}\b`)
	uriREg = regexp.MustCompile(`http[s]?\:\/\/.[a-zA-Z0-9\.\/\-]+`)
	newLineReg = regexp.MustCompile(`^[\r\n]+|\.|[\r\n]+$`)
)

// SentimentResult represents results of the job
type SentimentResult struct {
	ID        string    `json:"id"`
	QueryOn   time.Time `json:"ts"`
	Tweets    int     `json:"tweets"`
	NonRT    int     `json:"nonRT"`
	Sentiment int       `json:"sentiment"`
	Score     float32    `json:"score"`
	Magnitude float32    `json:"magnitude"`
}

func search(ctx context.Context, query string) (r *SentimentResult, err error) {

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		return nil, errors.New("Both, consumer key/secret and access token/secret are required")
	}

	// init convif
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// HTTP Client - will automatically authorize Requests
	httpClient := config.Client(ctx, token)
	client := twitter.NewClient(httpClient)

	searchArgs := &twitter.SearchTweetParams{
		Query:      query,
		Count:      100,
		Lang:       "en",
		ResultType: "recent",
	}

	log.Printf("Search: %v", query)
	search, _, err := client.Search.Tweets(searchArgs)
	if err != nil {
		return nil, err
	}

	// results
	result := &SentimentResult{
		ID:      newID(),
		Tweets:  len(search.Statuses),
		QueryOn: time.Now(),
	}

	contents := make([]string, 0)
	log.Printf("Found: %d", result.Tweets)
	for _, tweet := range search.Statuses {
		if tweet.RetweetedStatus == nil {
			//log.Printf("Raw: %s", tweet.Text)
			contents = append(contents, tweet.Text)
		}
	}

	// join
	result.NonRT = len(contents)

	// cleanup
	txt := strings.Join(contents, ". ")
	txt = userReg.ReplaceAllString(txt, " ")
	txt = uriREg.ReplaceAllString(txt, " ")
	txt = nonCharReg.ReplaceAllString(txt, " ")
	txt = shortReg.ReplaceAllString(txt, " ")
	txt = newLineReg.ReplaceAllString(txt, " ")

	// log.Printf("Text: %s", txt)
	sentiment, err := scoreSentiment(ctx, txt)
	if err != nil {
		log.Printf("Error while scoring: %v", err)
		return result, nil
	}

	log.Printf("Score: %f, Magnitude: %f", result.Score,sentiment.Magnitude)

	result.Magnitude = sentiment.Magnitude
	result.Score = sentiment.Score

	/*
		Clearly Positive*	"score": 0.8, 	"magnitude": 3.0
		Clearly Negative*	"score": -0.6, 	"magnitude": 4.0
		Neutral				"score": 0.1, 	"magnitude": 0.0
		Mixed				"score": 0.0, 	"magnitude": 4.0
	*/

	if result.Score > 0 && sentiment.Magnitude > magnitudeThreshold {
		result.Sentiment = 1
	} else if result.Score < 0 && sentiment.Magnitude > magnitudeThreshold {
		result.Sentiment = -1
	} else {
		result.Sentiment = 0
	}

	return result, nil

}

func newID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("Error while getting id: %v\n", err)
	}
	return fmt.Sprintf("qid-%s", id.String())
}
