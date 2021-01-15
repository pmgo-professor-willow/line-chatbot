package functions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/thoas/go-funk"
)

// RaidBoss is pre-processing data from LeekDuck website
type RaidBoss struct {
	Tier           string   `json:"tier"`
	No             int      `json:"no"`
	Name           string   `json:"name"`
	ImageURL       string   `json:"imageUrl"`
	ShinyAvailable bool     `json:"shinyAvailable"`
	Types          []string `json:"types"`
	TypeURLs       []string `json:"typeUrls"`
	CP             struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"cp"`
	BoostedCP struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"boostedCp"`
	BoostedWeathers    []string `json:"boostedWeathers"`
	BoostedWeatherURLs []string `json:"boostedWeatherUrls"`
}

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

// TweetMediaData is pre-processing data from Twitter API
type TweetMediaData struct {
	Type     string `json:"type"`
	URL      string `json:"url"`
	MediaKey string `json:"media_key"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

// TweetData is pre-processing data from Twitter API
type TweetData struct {
	ID        string         `json:"id"`
	Text      string         `json:"text"`
	Media     TweetMediaData `json:"media"`
	CreatedAt string         `json:"createdAt"`
}

// UserTweets is pre-processing data from Twitter API
type UserTweets struct {
	Name   string      `json:"name"`
	Tweets []TweetData `json:"tweets"`
}

// DataCache has all remote data and last updated time
type DataCache struct {
	RaidBosses []RaidBoss
	GameEvents []GameEvent
	TweetList  []UserTweets
	UpdatedAt  time.Time
}

// LoadRaidBosses load data from remote JSON
func LoadRaidBosses() []RaidBoss {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-leekduck/raid-bosses.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			events := []RaidBoss{}
			json.Unmarshal(bodyBuf, &events)

			return events
		}
	}

	return []RaidBoss{}
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

// LoadTweetList load data from remote JSON
func LoadTweetList() []UserTweets {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-tweets/tweet-list.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			tweetList := []UserTweets{}
			json.Unmarshal(bodyBuf, &tweetList)

			return tweetList
		}
	}

	return []UserTweets{}
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
		RaidBosses: []RaidBoss{},
		GameEvents: []GameEvent{},
		TweetList:  []UserTweets{},
		UpdatedAt:  time.Now().AddDate(0, 0, -1),
	}
}
