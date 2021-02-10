package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/thoas/go-funk"
)

// Event is pre-processing data from LeekDuck website
type Event struct {
	Title        string `json:"title"`
	Link         string `json:"link"`
	Type         string `json:"type"`
	ImageURL     string `json:"imageUrl"`
	Label        string `json:"label"`
	IsLocaleTime bool   `json:"isLocaleTime"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
}

// LoadEvents load data from remote JSON
func LoadEvents() []Event {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-leekduck/events.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			events := []Event{}
			json.Unmarshal(bodyBuf, &events)

			return events
		}
	}

	return []Event{}
}

// FilterEvents filters game events by specified label
func FilterEvents(gameEvents []Event, label string) []Event {
	loc, _ := time.LoadLocation(os.Getenv("TIMEZONE_LOCATION"))

	return funk.Filter(gameEvents, func(gameEvent Event) bool {
		isSameLabel := gameEvent.Label == label
		isMatchedEvent := false

		if gameEvent.Label == "current" {
			isInProgress := false
			if gameEvent.StartTime != "" && gameEvent.EndTime != "" {
				var startTime time.Time
				var endTime time.Time
				if gameEvent.IsLocaleTime {
					startTime, _ = time.ParseInLocation("2006-01-02T15:04:05Z", gameEvent.StartTime, loc)
					endTime, _ = time.ParseInLocation("2006-01-02T15:04:05Z", gameEvent.EndTime, loc)
				} else {
					startTime, _ = time.Parse(time.RFC3339, gameEvent.StartTime)
					endTime, _ = time.Parse(time.RFC3339, gameEvent.EndTime)
				}
				isInProgress = int(time.Now().Sub(startTime).Minutes()) > 0 && int(endTime.Sub(time.Now()).Minutes()) > 0
			} else if gameEvent.EndTime != "" {
				var endTime time.Time
				if gameEvent.IsLocaleTime {
					endTime, _ = time.ParseInLocation("2006-01-02T15:04:05Z", gameEvent.EndTime, loc)
				} else {
					endTime, _ = time.Parse(time.RFC3339, gameEvent.EndTime)
				}
				isInProgress = int(endTime.Sub(time.Now()).Minutes()) > 0
			} else if gameEvent.StartTime == "" && gameEvent.EndTime == "" {
				isInProgress = true
			}

			isMatchedEvent = isSameLabel && isInProgress
		} else if gameEvent.Label == "upcoming" {
			isWaiting := false
			if gameEvent.StartTime != "" {
				var startTime time.Time
				if gameEvent.IsLocaleTime {
					startTime, _ = time.ParseInLocation("2006-01-02T15:04:05Z", gameEvent.StartTime, loc)
				} else {
					startTime, _ = time.Parse(time.RFC3339, gameEvent.StartTime)
				}
				isWaiting = int(time.Now().Sub(startTime).Minutes()) < 0
			}

			isMatchedEvent = isSameLabel && isWaiting
		}

		return isMatchedEvent
	}).([]Event)
}
