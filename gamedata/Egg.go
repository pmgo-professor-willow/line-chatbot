package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/thoas/go-funk"
)

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
		ImageURL:       "https://sample.com/dummy.png",
		ShinyAvailable: false,
		CP: struct {
			Min int `json:"min"`
			Max int `json:"max"`
		}{
			Max: 0,
			Min: 0,
		},
	}
}
