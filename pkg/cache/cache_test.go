package cache

import (
	"testing"
	"time"

    "github.com/mchmarny/twfeel/pkg/common"

	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	v := toKey("knative")
	assert.Equal(t, v, "twfeel::a25hdGl2ZQ==")
}

func TestCache(t *testing.T) {

	r1 := &common.SentimentResult{
		ID: common.NewID(),
		Magnitude: 0.12,
		Query: "test",
		QueryOn: time.Now(),
		Score: 1.23,
		Tweets: 100,
	}

	err := Clear(r1.ID)
	assert.NoErrorf(t, err, "Error on Clear: %v", err)

	err = Set(r1.ID, r1)
	assert.NoErrorf(t, err, "Error on Set: %v", err)

	r2, err := Get(r1.ID)
	assert.NoErrorf(t, err, "Error on Get: %v", err)
	assert.NotNil(t, r2)
	assert.Equal(t, r1.ID, r2.ID)

}
