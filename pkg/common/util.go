package common

import (
	"fmt"
	"log"

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
