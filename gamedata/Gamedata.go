package gamedata

import (
	"time"
)

// DataCache has all remote data and last updated time
type DataCache struct {
	RaidBosses      []RaidBoss
	Eggs            []Egg
	Researches      []Research
	RocketInvasions []RocketInvasion
	Events          []Event
	TweetList       []UserTweets
	Channels        []Channel
	UpdatedAt       time.Time
}

// GetCache stores data from remote
func GetCache() *DataCache {
	return &DataCache{
		RaidBosses:      []RaidBoss{},
		Eggs:            []Egg{},
		Researches:      []Research{},
		RocketInvasions: []RocketInvasion{},
		Events:          []Event{},
		TweetList:       []UserTweets{},
		Channels:        []Channel{},
		UpdatedAt:       time.Now().AddDate(0, 0, -1),
	}
}
