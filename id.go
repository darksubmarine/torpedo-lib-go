package torpedo_lib

import (
	"github.com/google/uuid"
	ulid "github.com/oklog/ulid/v2"
)

// Ulid generates a ULID string. See https://github.com/ulid/spec
func Ulid() string {
	return ulid.Make().String()
}

// Uuid generates and inspects UUIDs based on RFC 4122
func Uuid() string {
	return uuid.NewString()
}
