package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/thoas/go-funk"
)

// Video is pre-processing data from YouTube API
type Video struct {
	Title        string `json:"title"`
	URL          string `json:"url"`
	Description  string `json:"description"`
	PublishedAt  string `json:"publishedAt"`
	ThumbnailURL string `json:"thumbnailUrl"`
	ChannelTitle string `json:"channelTitle"`
}

// Channel is pre-processing data from YouTube API
type Channel struct {
	Name            string  `json:"title"`
	ThumbnailURL    string  `json:"thumbnailUrl"`
	ViewCount       int     `json:"viewCount"`
	SubscriberCount int     `json:"subscriberCount"`
	Videos          []Video `json:"videos"`
}

// LoadChannels load data from remote JSON
func LoadChannels(cacheData *[]Channel) {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-youtuber/channels.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			channels := []Channel{}
			json.Unmarshal(bodyBuf, &channels)
			*cacheData = channels
		}
	}
}

// FindChannel finds channel by specified channel name
func FindChannel(channels []Channel, channelName string) Channel {
	return funk.Find(channels, func(channel Channel) bool {
		return channel.Name == channelName
	}).(Channel)
}
