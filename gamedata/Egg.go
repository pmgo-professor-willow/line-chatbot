package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/thoas/go-funk"
)

// Egg is pre-processing data from LeekDuck website
type Egg struct {
	No             int     `json:"no"`
	Name           string  `json:"name"`
	OriginalName   string  `json:"originalName"`
	Category       string  `json:"category"`
	ImageURL       string  `json:"imageUrl"`
	ShinyAvailable bool    `json:"shinyAvailable"`
	Regional       bool    `json:"regional"`
	CP             struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"cp"`
	Rate           float32 `json:"rate"`
}

// LoadEggs load data from remote JSON
func LoadEggs(cacheData *[]Egg) {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-thesilphroad/eggs.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			eggs := []Egg{}
			json.Unmarshal(bodyBuf, &eggs)
			*cacheData = eggs
		}
	}
}

// FilterdEggs filters eggs by specified category
func FilterdEggs(eggs []Egg, eggCategory string) []Egg {
	return funk.Filter(eggs, func(egg Egg) bool {
		return egg.Category == eggCategory
	}).([]Egg)
}

func CreateDummyEgg() Egg {
	return Egg{
		No:             0,
		Name:           " ",
		OriginalName:   "",
		Category:       "",
		ImageURL:       "https://example.com/dummy.png",
		ShinyAvailable: false,
		Regional:       false,
		CP: struct {
			Min int `json:"min"`
			Max int `json:"max"`
		}{
			Max: 0,
			Min: 0,
		},
		Rate:           0,
	}
}
