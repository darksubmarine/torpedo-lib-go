package tql

import (
	"fmt"
	"strings"
)

func sanitizeQryStr(s string) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(s, "'", "\\'"))
}

func sanitizeQryLikeStr(s string) string {
	return strings.ReplaceAll(s, "'", "\\'")
}
