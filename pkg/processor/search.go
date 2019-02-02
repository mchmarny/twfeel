package processor

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"
	"regexp"
	"strconv"

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

	// cache
	cacheMin = os.Getenv("CACHE_TTL_MIN")
	defaultCacheDuration = time.Minute * 5

	// validation expressions
	userReg = regexp.MustCompile(`@[\w]*`)
	nonCharReg = regexp.MustCompile(`[^a-zA-Z#]`)
	shortReg = regexp.MustCompile(`\b[a-z]{1,2}\b`)
	uriREg = regexp.MustCompile(`http[s]?\:\/\/.[a-zA-Z0-9\.\/\-]+`)
	newLineReg = regexp.MustCompile(`^[\r\n]+|\.|[\r\n]+$`)
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
		return nil, err
	}

	if cr != nil {
		log.Printf("Cache hit on `%s`", query)
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
	result := &common.SentimentResult{
		ID:      common.NewID(),
		Query:   query,
		QueryOn: time.Now(),
		Tweets:  len(search.Statuses),
	}


	contents := make([]string, 0)
	log.Printf("Found: %d", result.Tweets)
	for _, tweet := range search.Statuses {
		if tweet.RetweetedStatus == nil {
			//log.Printf("Raw: %s", tweet.Text)
			contents = append(contents, tweet.Text)
		}
	}

	// update after filter
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

	// set the cache
	resultTLL := defaultCacheDuration
	if cacheMin != "" {
		i, err := strconv.Atoi(cacheMin)
		if err != nil {
			log.Printf("CACHE_TTL_MIN set to invalid value (must be int): %v", err)
			return nil, err
		}
		resultTLL = time.Minute * time.Duration(i)
		log.Printf("Setting cache to %v based on CACHE_TTL_MIN", resultTLL)
	}

	err = cache.Set(query, result, resultTLL)
	if err != nil {
		log.Fatalf("BUG: Error while setting cache: %v", err)
	}

	return result, nil

}