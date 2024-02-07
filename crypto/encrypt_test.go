package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeString(t *testing.T) {
	cypherKey := []byte("asdfghjklzxcvbnmqwertyuiop123456")
	val := `This message is supper secret and must not be shared with no one else.
			Your mission should you choose to accept it...`

	encodedVal, err := EncodeString(cypherKey, val)
	assert.Nil(t, err)

	decodedVal, err := DecodeString(cypherKey, encodedVal)
	assert.Nil(t, err)
	assert.Equal(t, val, decodedVal)
}
