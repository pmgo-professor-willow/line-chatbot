package gamedata

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/thoas/go-funk"
)

// TweetMediaData is pre-processing data from Twitter API
type TweetMediaData struct {
	Type     string `json:"type"`
	URL      string `json:"url"`
	MediaKey string `json:"media_key"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

// TweetData is pre-processing data from Twitter API
type TweetData struct {
	ID        string           `json:"id"`
	Text      string           `json:"text"`
	MediaList []TweetMediaData `json:"mediaList"`
	CreatedAt string           `json:"createdAt"`
}

// UserTweets is pre-processing data from Twitter API
type UserTweets struct {
	Name   string      `json:"name"`
	Tweets []TweetData `json:"tweets"`
}

// LoadTweetList load data from remote JSON
func LoadTweetList(cacheData *[]UserTweets) {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-tweets/tweet-list.min.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			tweetList := []UserTweets{}
			json.Unmarshal(bodyBuf, &tweetList)
			*cacheData = tweetList
		}
	}
}

// FindUserTweets finds user tweets by specified user
func FindUserTweets(tweetList []UserTweets, twitterUser string) UserTweets {
	return funk.Find(tweetList, func(userTweets UserTweets) bool {
		return userTweets.Name == twitterUser
	}).(UserTweets)
}

// FindTweet finds tweet by specified ID
func FindTweet(tweets []TweetData, tweetID string) TweetData {
	return funk.Find(tweets, func(tweet TweetData) bool {
		return tweet.ID == tweetID
	}).(TweetData)
}
