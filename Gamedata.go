package functions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/thoas/go-funk"
)

// GameEvent is pre-processing data from LeekDuck website
type GameEvent struct {
	Title        string `json:"title"`
	Link         string `json:"link"`
	Type         string `json:"type"`
	ImageURL     string `json:"imageUrl"`
	Label        string `json:"label"`
	IsLocaleTime bool   `json:"isLocaleTime"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
}

// DataCache has all remote data and last updated time
type DataCache struct {
	GameEvents []GameEvent
	UpdatedAt  time.Time
}

// LoadGameEvents load data from remote JSON
func LoadGameEvents() []GameEvent {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-leekduck/events.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			events := []GameEvent{}
			json.Unmarshal(bodyBuf, &events)

			return events
		}
	}

	return []GameEvent{}
}

// FilterGameEvents filters game events by specified label
func FilterGameEvents(gameEvents []GameEvent, label string) []GameEvent {
	return funk.Filter(gameEvents, func(gameEvent GameEvent) bool {
		isCurrentEvent := gameEvent.Label == label

		isInProgress := false
		if gameEvent.StartTime != "" && gameEvent.EndTime != "" {
			startTime, _ := time.Parse(time.RFC3339, gameEvent.StartTime)
			endTime, _ := time.Parse(time.RFC3339, gameEvent.EndTime)
			isInProgress = int(time.Now().Sub(startTime).Minutes()) > 0 && int(endTime.Sub(time.Now()).Minutes()) > 0
		} else if gameEvent.EndTime != "" {
			endTime, _ := time.Parse(time.RFC3339, gameEvent.EndTime)
			isInProgress = int(endTime.Sub(time.Now()).Minutes()) > 0
		} else if gameEvent.StartTime == "" && gameEvent.EndTime == "" {
			isInProgress = true
		}

		return isCurrentEvent && isInProgress
	}).([]GameEvent)
}

// GetCache stores data from remote
func GetCache() *DataCache {
	return &DataCache{
		GameEvents: []GameEvent{},
		UpdatedAt:  time.Now().AddDate(0, 0, -1),
	}
}
