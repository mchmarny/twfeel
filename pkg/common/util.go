package common

import (
	"fmt"
	"log"
	"os"

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

// IsValidAccessToken centralizes token validation
func IsValidAccessToken(token string) bool {

	// set during tests
	if os.Getenv("SKIP_TW_TOKEN_VALIDATION") == "yes" {
		return true
	}
	// TODO: Allow for multiple known tokens
	return token == os.Getenv("ACCESS_TOKEN")
}
