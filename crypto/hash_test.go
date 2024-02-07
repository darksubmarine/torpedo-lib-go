package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHash(t *testing.T) {
	expected := "430ce34d020724ed75a196dfc2ad67c77772d169"

	hash := Hash("hello world!")
	assert.Equal(t, expected, hash)
}
