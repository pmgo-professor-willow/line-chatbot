package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"time"

	timeUtils "pmgo-professor-willow/lineChatbot/utils"

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
func LoadEvents(cacheData *[]Event) {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-leekduck/events.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			events := []Event{}
			json.Unmarshal(bodyBuf, &events)
			*cacheData = events
		}
	}
}

// FilterEvents filters game events by specified label
func FilterEvents(gameEvents []Event, eventLabel string, eventType interface{}) []Event {
	filteredEvents := funk.Filter(gameEvents, func(gameEvent Event) bool {
		isSameLabel := gameEvent.Label == eventLabel
		isMatchedEvent := false

		if eventType != nil && eventType != "" && gameEvent.Type != eventType {
			return false
		}

		if gameEvent.Label == "current" {
			isInProgress := false
			if gameEvent.StartTime != "" && gameEvent.EndTime != "" {
				startTime := timeUtils.ToTimeInstance(gameEvent.StartTime, gameEvent.IsLocaleTime)
				endTime := timeUtils.ToTimeInstance(gameEvent.EndTime, gameEvent.IsLocaleTime)
				isInProgress = int(time.Now().Sub(startTime).Minutes()) > 0 && int(endTime.Sub(time.Now()).Minutes()) > 0
			} else if gameEvent.EndTime != "" {
				endTime := timeUtils.ToTimeInstance(gameEvent.EndTime, gameEvent.IsLocaleTime)
				isInProgress = int(endTime.Sub(time.Now()).Minutes()) > 0
			} else if gameEvent.StartTime == "" && gameEvent.EndTime == "" {
				isInProgress = true
			}

			isMatchedEvent = isSameLabel && isInProgress
		} else if gameEvent.Label == "upcoming" {
			isWaiting := false
			if gameEvent.StartTime != "" {
				startTime := timeUtils.ToTimeInstance(gameEvent.StartTime, gameEvent.IsLocaleTime)
				isWaiting = int(time.Now().Sub(startTime).Minutes()) < 0
			}

			isMatchedEvent = isSameLabel && isWaiting
		}

		return isMatchedEvent
	}).([]Event)

	if eventLabel == "current" {
		sort.SliceStable(filteredEvents, func(i, j int) bool {
			return filteredEvents[i].EndTime < filteredEvents[j].EndTime
		})
	} else if eventLabel == "upcoming" {
		sort.SliceStable(filteredEvents, func(i, j int) bool {
			return filteredEvents[i].StartTime < filteredEvents[j].StartTime
		})
	}

	return filteredEvents
}
