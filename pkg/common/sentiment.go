package common

import (
	"time"
)

// SentimentResult represents results of the job
type SentimentResult struct {
	ID        string    `json:"id"`
	Query     string    `json:"query"`
	QueryOn   time.Time `json:"ts"`
	Tweets    int       `json:"tweets"`
	NonRT     int       `json:"nonRT"`
	Score     float32   `json:"score"`
	Magnitude float32   `json:"magnitude"`
}
