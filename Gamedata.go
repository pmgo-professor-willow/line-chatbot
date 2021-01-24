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
	OriginalName   string   `json:"originalName"`
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

// Egg is pre-processing data from LeekDuck website
type Egg struct {
	No             int    `json:"no"`
	Name           string `json:"name"`
	OriginalName   string `json:"originalName"`
	Category       string `json:"category"`
	ImageURL       string `json:"imageUrl"`
	ShinyAvailable bool   `json:"shinyAvailable"`
	CP             struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"cp"`
}

// ResearchRewardPokemon is pre-processing data from The Silph Road website
type ResearchRewardPokemon struct {
	No             int    `json:"no"`
	Name           string `json:"name"`
	OriginalName   string `json:"originalName"`
	ShinyAvailable bool   `json:"shinyAvailable"`
	ImageURL       string `json:"imageUrl"`
	CP             struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"cp"`
}

// Research is pre-processing data from The Silph Road website
type Research struct {
	Description         string                  `json:"description"`
	OriginalDescription string                  `json:"originalDescription"`
	Category            string                  `json:"category"`
	RewardPokemons      []ResearchRewardPokemon `json:"rewardPokemons"`
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

// Video is pre-processing data from YouTube API
type Video struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	Description  string `json:"description"`
	PublishedAt  string `json:"publishedAt"`
	ThumbnailURL string `json:"thumbnailUrl"`
	ChannelTitle string `json:"channelTitle"`
}

// Channel is pre-processing data from YouTube API
type Channel struct {
	Name            string  `json:"title"`
	ThumbnailURL    string  `json:"thumbnailUrl"`
	ViewCount       int     `json:"viewCount"`
	SubscriberCount int     `json:"subscriberCount"`
	Videos          []Video `json:"videos"`
}

// DataCache has all remote data and last updated time
type DataCache struct {
	RaidBosses []RaidBoss
	Eggs       []Egg
	Researches []Research
	GameEvents []GameEvent
	TweetList  []UserTweets
	Channels   []Channel
	UpdatedAt  time.Time
}

// LoadRaidBosses load data from remote JSON
func LoadRaidBosses() []RaidBoss {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-leekduck/raid-bosses.min.json")

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

// LoadEggs load data from remote JSON
func LoadEggs() []Egg {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-leekduck/eggs.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			eggs := []Egg{}
			json.Unmarshal(bodyBuf, &eggs)

			return eggs
		}
	}

	return []Egg{}
}

// LoadResearches load data from remote JSON
func LoadResearches() []Research {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-thesilphroad/researches.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			researchs := []Research{}
			json.Unmarshal(bodyBuf, &researchs)

			return researchs
		}
	}

	return []Research{}
}

// LoadGameEvents load data from remote JSON
func LoadGameEvents() []GameEvent {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-leekduck/events.min.json")

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
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-tweets/tweet-list.min.json")

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

// LoadChannels load data from remote JSON
func LoadChannels() []Channel {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-youtuber/channels.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			channels := []Channel{}
			json.Unmarshal(bodyBuf, &channels)

			return channels
		}
	}

	return []Channel{}
}

// FilterdRaidBosses filters raid bosses by specified tier
func FilterdRaidBosses(raidBosses []RaidBoss, raidTier string) []RaidBoss {
	return funk.Filter(raidBosses, func(raidBoss RaidBoss) bool {
		return raidBoss.Tier == raidTier
	}).([]RaidBoss)
}

// FilterdEggs filters eggs by specified category
func FilterdEggs(eggs []Egg, eggCategory string) []Egg {
	return funk.Filter(eggs, func(egg Egg) bool {
		return egg.Category == eggCategory
	}).([]Egg)
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

// FindUserTweets finds user tweets by specified user
func FindUserTweets(tweetList []UserTweets, twitterUser string) UserTweets {
	return funk.Find(tweetList, func(userTweets UserTweets) bool {
		return userTweets.Name == twitterUser
	}).(UserTweets)
}

// FindTweet finds tweet by specified ID
func FindTweet(tweets []TweetData, tweetID string) TweetData {
	return funk.Find(tweets, func(tweet TweetData) bool {
		return tweet.ID == tweetID
	}).(TweetData)
}

// FindChannel finds channel by specified channel name
func FindChannel(channels []Channel, channelName string) Channel {
	return funk.Find(channels, func(channel Channel) bool {
		return channel.Name == channelName
	}).(Channel)
}

// GetCache stores data from remote
func GetCache() *DataCache {
	return &DataCache{
		RaidBosses: []RaidBoss{},
		Eggs:       []Egg{},
		Researches: []Research{},
		GameEvents: []GameEvent{},
		TweetList:  []UserTweets{},
		Channels:   []Channel{},
		UpdatedAt:  time.Now().AddDate(0, 0, -1),
	}
}
