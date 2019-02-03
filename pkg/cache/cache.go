package cache

import (
	"fmt"
	"log"
	"os"
	"time"
	"strconv"

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

	defaultCacheDurationMin = 5
	cacheDuration = time.Minute * time.Duration(defaultCacheDurationMin)

)

func init() {

	cacheMin := os.Getenv("CACHE_TTL_MIN")
	if cacheMin != "" {
		i, err := strconv.Atoi(cacheMin)
		if err != nil {
			log.Printf(
				"Error: CACHE_TTL_MIN set to invalid value (want int): %v. Using default (%d min)",
				err, defaultCacheDurationMin)
		} else {
			cacheDuration = time.Minute * time.Duration(i)
			log.Printf("Cache set to %v based on CACHE_TTL_MIN", cacheDuration)
		}
	}

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

// Clear clears cache for specific key
func Clear(key string) error {

	if key == "" {
		return fmt.Errorf("Nil key")
	}

	k := toKey(key)

	if !codec.Exists(k){
		return nil
	}

	return codec.Delete(toKey(key))

}

// Set sets the result cache
func Set(key string, rez *common.SentimentResult) error {

	if key == "" {
		return fmt.Errorf("Nil key")
	}

	if rez == nil {
		return fmt.Errorf("Nil result to cache")
	}

	return codec.Set(&cache.Item{
		Key:        toKey(key),
		Object:     rez,
		Expiration: cacheDuration,
	})
}

// Get retreaves result from cache
func Get(key string) (rez *common.SentimentResult, err error) {

	if key == "" {
		return nil, fmt.Errorf("Nil key")
	}

	k := toKey(key)

	if !codec.Exists(k){
		return nil, nil
	}

	var sr common.SentimentResult
	e := codec.Get(k, &sr)

	if e != nil {
		return nil, e
	}

	if sr.ID == "" {
		return nil, nil
	}

	return &sr, nil

}
