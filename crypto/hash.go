package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func Hash(str string) string {
	return hex.EncodeToString(md5.New().Sum([]byte(str)))
}
