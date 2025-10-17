package helper

import (
	"time"
)

var startTime = time.Now()

func Hello() int {
	uptimeMinutes := int(time.Since(startTime).Minutes())
	return uptimeMinutes
}
