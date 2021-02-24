package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

// LoadRaidBosses load data from remote JSON
func LoadRaidBosses(cacheData *[]RaidBoss) {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-leekduck/raid-bosses.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			raidBosses := []RaidBoss{}
			json.Unmarshal(bodyBuf, &raidBosses)
			*cacheData = raidBosses
		}
	}
}

// FilterdRaidBosses filters raid bosses by specified tier
func FilterdRaidBosses(raidBosses []RaidBoss, raidTier string) []RaidBoss {
	return funk.Filter(raidBosses, func(raidBoss RaidBoss) bool {
		return raidBoss.Tier == raidTier
	}).([]RaidBoss)
}
