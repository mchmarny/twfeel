package common

import (
	"time"
)

const (
	magnitudeThreshold = 0.2

	// UndefinedSentiment is not set
	UndefinedSentiment SentimentMeasure = 0
	// PositiveSentiment is good
	PositiveSentiment SentimentMeasure = 1
	// NegativeSentiment is bad
	NegativeSentiment SentimentMeasure = 3
	// MixedSentiment is meh
	MixedSentiment SentimentMeasure = 4
)

// SentimentMeasure defines the feeling on topic
type SentimentMeasure int

func (s SentimentMeasure) String() string {
	names := [...]string{
		"UndefinedSentiment",
		"PositiveSentiment",
		"NegativeSentiment",
		"MixedSentiment",
	}

	if s < PositiveSentiment || s > MixedSentiment {
		return "UndefinedSentiment"
	}

	return names[s]
}

// IsPositive is a shortcut
func (s SentimentMeasure) IsPositive() bool {
	return s == PositiveSentiment
}

// IsNegative is a shortcut
func (s SentimentMeasure) IsNegative() bool {
	return s == NegativeSentiment
}

// SentimentResult represents results of the job
type SentimentResult struct {
	ID        string           `json:"id"`
	Query     string           `json:"query"`
	QueryOn   time.Time        `json:"ts"`
	Tweets    int              `json:"tweets"`
	Score     float32          `json:"score"`
	Magnitude float32          `json:"magnitude"`
	Sentiment SentimentMeasure `json:"sentiment"`
}

// CalculateSentiment derives sentiment
func (s *SentimentResult) CalculateSentiment() SentimentMeasure {
	if s.Score > 0.2 && s.Magnitude > magnitudeThreshold {
		s.Sentiment = PositiveSentiment
	} else if s.Score < -0.2 && s.Magnitude > magnitudeThreshold {
		s.Sentiment = NegativeSentiment
	} else {
		s.Sentiment = MixedSentiment
	}
	return s.Sentiment
}

/*
	Clearly Positive*	"score": 0.8, 	"magnitude": 3.0
	Clearly Negative*	"score": -0.6, 	"magnitude": 4.0
	Neutral				"score": 0.1, 	"magnitude": 0.0
	Mixed				"score": 0.0, 	"magnitude": 4.0
*/
