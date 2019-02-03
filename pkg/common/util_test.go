package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidAccessToken(t *testing.T) {
	tokens := "rest:1111;chat:2222;slack:3333"
	assert.True(t, isValidToken(tokens, "rest", "1111"))
	assert.True(t, isValidToken(tokens, "chat", "2222"))
	assert.True(t, isValidToken(tokens, "slack", "3333"))
}
