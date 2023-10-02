package torpedo_lib

import "time"

func TimeNow() int64 {
	return time.Now().UnixMilli()
}
