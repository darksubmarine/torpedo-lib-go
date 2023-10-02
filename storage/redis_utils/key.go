package redis_utils

import (
	"github.com/darksubmarine/torpedo-lib-go"
)

func EntityKey(name, id string) string {
	return torpedo_lib.NewTRN("entity", name, id).String()
}
