package gamedata

import (
	"time"
)

// DataCache has all remote data and last updated time
type DataCache struct {
	RaidBosses               []RaidBoss
	RaidBossesUpdatedAt      time.Time
	Eggs                     []Egg
	EggsUpdatedAt            time.Time
	Researches               []Research
	ResearchesUpdatedAt      time.Time
	RocketInvasions          []RocketInvasion
	RocketInvasionsUpdatedAt time.Time
	Events                   []Event
	EventsUpdatedAt          time.Time
	TweetList                []UserTweets
	TweetListUpdatedAt       time.Time
	Channels                 []Channel
	ChannelsUpdatedAt        time.Time
}

// GetCache stores data from remote
func GetCache() *DataCache {
	return &DataCache{
		RaidBosses:               []RaidBoss{},
		RaidBossesUpdatedAt:      time.Now().AddDate(0, 0, -1),
		Eggs:                     []Egg{},
		EggsUpdatedAt:            time.Now().AddDate(0, 0, -1),
		Researches:               []Research{},
		ResearchesUpdatedAt:      time.Now().AddDate(0, 0, -1),
		RocketInvasions:          []RocketInvasion{},
		RocketInvasionsUpdatedAt: time.Now().AddDate(0, 0, -1),
		Events:                   []Event{},
		EventsUpdatedAt:          time.Now().AddDate(0, 0, -1),
		TweetList:                []UserTweets{},
		TweetListUpdatedAt:       time.Now().AddDate(0, 0, -1),
		Channels:                 []Channel{},
		ChannelsUpdatedAt:        time.Now().AddDate(0, 0, -1),
	}
}
