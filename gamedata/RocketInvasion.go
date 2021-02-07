package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/thoas/go-funk"
)

// LineupPokemon is pre-processing data from YouTube API
type LineupPokemon struct {
	SlotNo         int    `json:"slotNo"`
	No             int    `json:"no"`
	Name           string `json:"name"`
	OriginalName   string `json:"originalName"`
	Catchable      bool   `json:"catchable"`
	ShinyAvailable bool   `json:"shinyAvailable"`
	ImageURL       string `json:"imageUrl"`
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
func LoadRocketInvasions() []RocketInvasion {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-thesilphroad/rocket-invasions.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			rocketInvasions := []RocketInvasion{}
			json.Unmarshal(bodyBuf, &rocketInvasions)

			return rocketInvasions
		}
	}

	return []RocketInvasion{}
}

// FilterdRocketInvasions filters rocket invasions by specified label
func FilterdRocketInvasions(rocketInvasions []RocketInvasion, label string) []RocketInvasion {
	return funk.Filter(rocketInvasions, func(rocketInvasion RocketInvasion) bool {
		isGrunt := label == "grunt" && !rocketInvasion.IsSpecial
		isSpecial := label == "special" && rocketInvasion.IsSpecial
		return isGrunt || isSpecial
	}).([]RocketInvasion)
}
