package processor

import (
	"context"
	"errors"
	"log"
	"fmt"
	"os"
	"strings"
	"time"
	"regexp"

	"github.com/mchmarny/twfeel/pkg/common"
	"github.com/mchmarny/twfeel/pkg/cache"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

)

var (
	// twitter API secrets
	consumerKey = os.Getenv("T_CONSUMER_KEY")
	consumerSecret = os.Getenv("T_CONSUMER_SECRET")
	accessToken = os.Getenv("T_ACCESS_TOKEN")
	accessSecret = os.Getenv("T_ACCESS_SECRET")

	// validation expressions
	userReg = regexp.MustCompile(`@[\w]*`)
	nonWordReg = regexp.MustCompile(`[^a-zA-Z#]`)
	shortReg = regexp.MustCompile(`\b[a-z]{1,2}\b`)
	uriReg = regexp.MustCompile(`http[s]?\:\/\/.[a-zA-Z0-9\.\/\-]+`)
	emailReg = regexp.MustCompile(`\S*@\S*\s?`)
	spaceReg = regexp.MustCompile(`\s+`)
)

// Search searches Twitter and scores results
func Search(ctx context.Context, query string) (r *common.SentimentResult, err error) {

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		return nil, errors.New("Both, consumer key/secret and access token/secret are required")
	}

	// check cache
	cr, err := cache.Get(query)
	if err != nil {
		log.Printf("Error quering cache: %v", err)
	}

	if cr != nil && cr.IsValidState() {
		log.Printf("Cache hit on `%s`", query)
		// TODO: refactor this to allow cached objects to persist sentiment
		cr.CalculateSentiment()
		return cr, nil
	}

	// init convif
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// HTTP Client - will automatically authorize Requests
	httpClient := config.Client(ctx, token)
	client := twitter.NewClient(httpClient)


	// TODO: Add debugging for rate limits
	// rateArgs := &twitter.RateLimitParams{}
	// client.RateLimits.Status()

	searchArgs := &twitter.SearchTweetParams{
		Query:      fmt.Sprintf("%s -filter:retweets", query),
		Count:      100,
		Lang:       "en",
		ResultType: "recent",
	}

	log.Printf("Searching: %v", query)
	search, _, err := client.Search.Tweets(searchArgs)
	if err != nil {
		return nil, err
	}

	// results
	result := &common.SentimentResult{
		ID:      common.NewID(),
		Query:   query,
		QueryOn: time.Now(),
		Tweets:  len(search.Statuses),
	}

	contents := make([]string, 0)
	log.Printf("Found: %d", result.Tweets)
	for _, tweet := range search.Statuses {
		//log.Printf("Raw: %s", tweet.Text)
		contents = append(contents, tweet.Text)
	}

	// clean
	txt := toCleanString(contents)

	// log.Printf("Text: %s", txt)
	sentiment, err := scoreSentiment(ctx, txt)
	if err != nil {
		log.Printf("Error while scoring: %v", err)
		return result, nil
	}

	log.Printf("Score: %f, Magnitude: %f", result.Score, sentiment.Magnitude)
	result.Magnitude = sentiment.Magnitude
	result.Score = sentiment.Score

	// after the score and magnitude are set, calc sentiment
	result.CalculateSentiment()

	// set the cache
	err = cache.Set(query, result)
	if err != nil {
		log.Fatalf("BUG: Error while setting cache: %v", err)
	}

	return result, nil

}

func toCleanString(contents []string) string {

	txt := strings.Join(contents, ". ")

	txt = strings.ToLower(txt)

	// remove email addresses
	txt = emailReg.ReplaceAllString(txt, " ")

	// remove @usernames
	txt = userReg.ReplaceAllString(txt, " ")

	// remove URLs
	txt = uriReg.ReplaceAllString(txt, " ")

	// remove non-words
	txt = nonWordReg.ReplaceAllString(txt, " ")

	// remove short words (<3 char)
	txt = shortReg.ReplaceAllString(txt, " ")

	// space
	txt = spaceReg.ReplaceAllString(txt, " ")

	txt = strings.Trim(txt, " ")

	return txt

}