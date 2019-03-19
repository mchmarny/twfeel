package common

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

// NewID generates new globally unique ID
func NewID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("Error while getting id: %v\n", err)
	}
	return fmt.Sprintf("id-%s", id.String())
}

var tokenMap map[string]string

// IsValidAccessToken centralizes token validation
func IsValidAccessToken(key, token string) bool {

	// set during tests
	if os.Getenv("SKIP_TW_TOKEN_VALIDATION") == "yes" {
		return true
	}
	tokens := os.Getenv("ACCESS_TOKENS")

	if tokens == "" {
		log.Println("ACCESS_TOKENS not set")
		return false
	}

	return isValidToken(tokens, key, token)
}

func isValidToken(tokens, key, token string) bool {
	if tokenMap == nil {
		tokenMap = make(map[string]string)
		tokenParts := strings.Split(tokens, ";")
		for _, part := range tokenParts {
			tokenPairs := strings.Split(part, ":")
			if len(tokenPairs) == 2 {
				k := strings.Trim(tokenPairs[0], " ")
				v := strings.Trim(tokenPairs[1], " ")
				tokenMap[k] = v
			}
		}
	}
	return tokenMap[key] == token
}
