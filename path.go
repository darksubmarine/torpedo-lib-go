package torpedo_lib

import (
	"path"
)

func P(elem ...string) string {
	return path.Join(elem...)
}
