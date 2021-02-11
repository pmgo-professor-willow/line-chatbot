package utils

import (
	"os"
	"time"
)

var loc, _ = time.LoadLocation(os.Getenv("TIMEZONE_LOCATION"))

func ToTimeInstance(timeText string, isLocaleTime bool) time.Time {
	var formattedTime time.Time

	if isLocaleTime {
		formattedTime, _ = time.ParseInLocation("2006-01-02T15:04:05Z", timeText, loc)
	} else {
		formattedTime, _ = time.Parse(time.RFC3339, timeText)
	}

	return formattedTime
}
