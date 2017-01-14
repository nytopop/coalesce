// coalesce/util/time.go

package util

import (
	"strconv"
	"time"
)

func NiceTime(oldTime int64) string {
	curTime := time.Now().Unix()
	seconds := curTime - oldTime
	var elapsed string

	switch {
	// < 2 minutes
	case seconds < 120:
		elapsed = strconv.Itoa(int(seconds))
		return elapsed + " seconds ago"

	// < 2 hours
	case seconds < 7200:
		elapsed = strconv.Itoa(int(seconds / 60))
		return elapsed + " minutes ago"

	// < 2 days
	case seconds < 172800:
		elapsed = strconv.Itoa(int(seconds / 60 / 60))
		return elapsed + " hours ago"

	// < 2 months
	case seconds < 5256000:
		elapsed = strconv.Itoa(int(seconds / 60 / 60 / 24))
		return elapsed + " days ago"

	// < 2 years
	case seconds < 63072000:
		elapsed = strconv.Itoa(int(seconds / 60 / 60 / 24 / 30))
		return elapsed + " months ago"

	// 2 years +
	default:
		elapsed = strconv.Itoa(int(seconds / 60 / 60 / 24 / 30 / 12))
		return elapsed + " years ago"
	}
}
