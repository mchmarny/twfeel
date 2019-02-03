package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCleanString(t *testing.T) {

	dirtyStr := []string{
		`Test @username gr8 job
		 here is`,

		`Let's see how it goes k`,

		`My url is http://test.com, my email user@test.com`,
	}

	expected := "test job here let see how goes url email"

	got := toCleanString(dirtyStr)

	assert.Equal(t, expected, got)
}
