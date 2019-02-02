package cache

import (
	"fmt"
	"log"
	"os"
	"time"

	b64 "encoding/base64"

   "github.com/mchmarny/twfeel/pkg/common"

	"github.com/go-redis/redis"
	"github.com/go-redis/cache"
	"github.com/vmihailenco/msgpack"
)

const (
	cacheKeyPrefix = "twfeel"
	redisHostToken = "REDIS_HOST"
	redisPassToken = "REDIS_PASS"
)

var (
	codec *cache.Codec
)

func init() {

	host := os.Getenv(redisHostToken)
	if host == "" {
		log.Fatalf("Required variable undefined: %s", redisHostToken)
	}

	pass := os.Getenv(redisPassToken)
	if pass == "" {
		log.Fatalf("Required variable undefined: %s", redisPassToken)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass,
		DB:       0,
	})

	log.Printf("Connecting to %s...", host)
	pong, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Error on PING: %v", err)
	}

	log.Printf("Success, PING=%s", pong)

	codec = &cache.Codec{
		Redis: client,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}

}

func toKey(val string) string {
	return fmt.Sprintf("%s::%s", cacheKeyPrefix,
		b64.StdEncoding.EncodeToString([]byte(val)))
}

// Set sets the result cache
func Set(key string, rez *common.SentimentResult, ttl time.Duration) error {

	if key == "" {
		return fmt.Errorf("Nil key")
	}

	if rez == nil {
		return fmt.Errorf("Nil result to cache")
	}

	return codec.Set(&cache.Item{
		Key:        toKey(key),
		Object:     rez,
		Expiration: ttl,
	})
}

// Get retreaves result from cache
func Get(key string) (rez *common.SentimentResult, err error) {

	if key == "" {
		return nil, fmt.Errorf("Nil key")
	}

	var sr common.SentimentResult
	e := codec.Get(toKey(key), &sr)

	if e == redis.Nil {
		return nil, nil
	}

	return &sr, nil

}
