package helper

import (
	"time"
)

var startTime = time.Now()

func Uptime() int {
	uptimeMinutes := int(time.Since(startTime).Minutes())
	return uptimeMinutes
}
