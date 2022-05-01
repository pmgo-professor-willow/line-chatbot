package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/thoas/go-funk"
)

// ResearchRewardPokemon is pre-processing data from The Silph Road website
type ResearchRewardPokemon struct {
	No             int    `json:"no"`
	Name           string `json:"name"`
	OriginalName   string `json:"originalName"`
	ShinyAvailable bool   `json:"shinyAvailable"`
	CP             struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"cp"`
	ImageURL string `json:"imageUrl"`
}

// ResearchRewardPokemonMegaCandy is pre-processing data from The Silph Road website
type ResearchRewardPokemonMegaCandy struct {
	No                int    `json:"no"`
	Name              string `json:"name"`
	OriginalName      string `json:"originalName"`
	Count             int    `json:"count"`
	ImageURL          string `json:"imageUrl"`
	MegaCandyImageUrl string `json:"megaCandyImageUrl"`
}

// Research is pre-processing data from The Silph Road website
type Research struct {
	Description              string                           `json:"description"`
	OriginalDescription      string                           `json:"originalDescription"`
	Category                 string                           `json:"category"`
	RewardPokemons           []ResearchRewardPokemon          `json:"rewardPokemons"`
	RewardPokemonMegaCandies []ResearchRewardPokemonMegaCandy `json:"rewardPokemonMegaCandies"`
}

// LoadResearches load data from remote JSON
func LoadResearches(cacheData *[]Research) {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-thesilphroad/researches.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			researches := []Research{}
			json.Unmarshal(bodyBuf, &researches)
			*cacheData = researches
		}
	}
}

// FilterResearches filters researach by specified label
func FilterResearches(researches []Research, label string) []Research {
	return funk.Filter(researches, func(research Research) bool {
		isEvent := label == "event" && research.Category == "活動限定"
		isNotEvent := label == "others" && research.Category != "活動限定"

		return isEvent || isNotEvent
	}).([]Research)
}
