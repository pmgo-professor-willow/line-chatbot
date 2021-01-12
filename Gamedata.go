package functions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type GameEvent struct {
	Title        string `json:"title"`
	Link         string `json:"link"`
	Type         string `json:"type"`
	ImageUrl     string `json:"imageUrl"`
	Label        string `json:"label"`
	IsLocaleTime bool   `json:"isLocaleTime"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
}

type DataCache struct {
	GameEvents []GameEvent
	UpdatedAt  time.Time
}

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

func GetCache() *DataCache {
	return &DataCache{
		GameEvents: []GameEvent{},
		UpdatedAt:  time.Now().AddDate(0, 0, -1),
	}
}
