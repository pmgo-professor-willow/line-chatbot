package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/thoas/go-funk"
)

// LineupPokemon is pre-processing data from YouTube API
type LineupPokemon struct {
	SlotNo         int      `json:"slotNo"`
	No             int      `json:"no"`
	Name           string   `json:"name"`
	OriginalName   string   `json:"originalName"`
	Types          []string `json:"types"`
	Catchable      bool     `json:"catchable"`
	ShinyAvailable bool     `json:"shinyAvailable"`
	ImageURL       string   `json:"imageUrl"`
}

// RocketInvasion is pre-processing data from YouTube API
type RocketInvasion struct {
	Quote             string          `json:"quote"`
	OrignialQuote     string          `json:"orignialQuote"`
	Category          string          `json:"category"`
	CharacterImageURL string          `json:"characterImageUrl"`
	IsSpecial         bool            `json:"isSpecial"`
	LineupPokemons    []LineupPokemon `json:"lineupPokemons"`
}

// LoadRocketInvasions load data from remote JSON
func LoadRocketInvasions(cacheData *[]RocketInvasion) {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-pokemongohub/rocketInvasions.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			rocketInvasions := []RocketInvasion{}
			json.Unmarshal(bodyBuf, &rocketInvasions)
			*cacheData = rocketInvasions
		}
	}
}

// FilterRocketInvasions filters rocket invasions by specified label
func FilterRocketInvasions(rocketInvasions []RocketInvasion, label string) []RocketInvasion {
	return funk.Filter(rocketInvasions, func(rocketInvasion RocketInvasion) bool {
		isGrunt := label == "grunt" && !rocketInvasion.IsSpecial
		isSpecial := label == "special" && rocketInvasion.IsSpecial
		return isGrunt || isSpecial
	}).([]RocketInvasion)
}

// FindRocketInvasion finds rocket invasion by specified category
func FindRocketInvasion(rocketInvasions []RocketInvasion, category string) RocketInvasion {
	return funk.Find(rocketInvasions, func(rocketInvasion RocketInvasion) bool {
		return rocketInvasion.Category == category
	}).(RocketInvasion)
}
