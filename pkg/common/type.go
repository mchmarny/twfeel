package common

import (
	"time"
)

const (
	magnitudeThreshold = 0.2

	// UndefinedSentiment is not set
	UndefinedSentiment = "Undefined"
	// PositiveSentiment is good
	PositiveSentiment = "Positive"
	// NegativeSentiment is bad
	NegativeSentiment = "Negative"
	// MixedSentiment is meh
	MixedSentiment = "Mixed"
	// NeutralSentiment is just that
	NeutralSentiment = "Neutral"
)

// SentimentResult represents results of the job
type SentimentResult struct {
	ID        string    `json:"id"`
	Query     string    `json:"query"`
	QueryOn   time.Time `json:"ts"`
	Tweets    int       `json:"tweets"`
	Score     float32   `json:"score"`
	Magnitude float32   `json:"magnitude"`
	Sentiment string    `json:"sentiment"`
}

// IsValidState checks if the result was populated
func (s *SentimentResult) IsValidState() bool {
	return s.ID != "" && s.Query != "" && s.Sentiment != UndefinedSentiment
}

// CalculateSentiment derives sentiment
func (s *SentimentResult) CalculateSentiment() string {
	if s.Score > 0.2 && s.Magnitude > magnitudeThreshold {
		s.Sentiment = PositiveSentiment
	} else if s.Score < -0.2 && s.Magnitude > magnitudeThreshold {
		s.Sentiment = NegativeSentiment
	} else if (s.Score > -0.2 && s.Score < 0.2) && s.Magnitude > 2.0 {
		s.Sentiment = MixedSentiment
	} else {
		s.Sentiment = NeutralSentiment
	}
	return s.Sentiment
}

/*
	Clearly Positive*	"score": 0.8, 	"magnitude": 3.0
	Clearly Negative*	"score": -0.6, 	"magnitude": 4.0
	Neutral				"score": 0.1, 	"magnitude": 0.0
	Mixed				"score": 0.0, 	"magnitude": 4.0
*/
