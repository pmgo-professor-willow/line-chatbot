package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

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
