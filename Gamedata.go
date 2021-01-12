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

type TweetMediaData struct {
	Type     string `json:"type"`
	URL      string `json:"url"`
	MediaKey string `json:"media_key"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

type TweetData struct {
	Id        string         `json:"id"`
	Text      string         `json:"text"`
	Media     TweetMediaData `json:"media"`
	CreatedAt string         `json:"createdAt"`
}

type UserTweets struct {
	Name   string      `json:"name"`
	Tweets []TweetData `json:"tweets"`
}

type DataCache struct {
	GameEvents []GameEvent
	TweetList  []UserTweets
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

func GetCache() *DataCache {
	return &DataCache{
		GameEvents: []GameEvent{},
		TweetList:  []UserTweets{},
		UpdatedAt:  time.Now().AddDate(0, 0, -1),
	}
}
