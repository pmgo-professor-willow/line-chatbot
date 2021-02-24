package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// OfficialPost is pre-processing data from official website
type OfficialPost struct {
	Title         string `json:"title"`
	Link          string `json:"link"`
	Date          string `json:"date"`
	CoverImageURL string `json:"coverImageUrl"`
}

// LoadOfficialPosts load data from remote JSON
func LoadOfficialPosts(cacheData *[]OfficialPost) {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-pokemongolive/posts.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			officialPosts := []OfficialPost{}
			json.Unmarshal(bodyBuf, &officialPosts)
			*cacheData = officialPosts
		}
	}
}
