package functions

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

// WebhookFunction is base LINE webhook entry
func WebhookFunction(w http.ResponseWriter, r *http.Request) {
	cache := GetCache()

	// Refresh cache about data from cloud.
	if time.Since(cache.UpdatedAt).Minutes() > 1 {
		cache.RaidBosses = LoadRaidBosses()
		cache.GameEvents = LoadGameEvents()
		cache.TweetList = LoadTweetList()
		cache.UpdatedAt = time.Now()
	}

	// LINE messaging API client
	var client, _ = linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		linebot.WithHTTPClient(&http.Client{}),
	)

	events, _ := client.ParseRequest(r)
	for _, event := range events {
		if event.Type != linebot.EventTypePostback {
			break
		}

		qs, _ := url.ParseQuery(event.Postback.Data)

		if qs.Get("raidTier") != "" {
			selectedRaidTier := qs.Get("raidTier")
			selectedRaidBosses := FilterdRaidBosses(cache.RaidBosses, selectedRaidTier)
			messages := GenerateRaidBossMessages(selectedRaidBosses, selectedRaidTier)

			replyMessageCall := client.ReplyMessage(event.ReplyToken, messages...)

			if _, err := replyMessageCall.Do(); err != nil {
			}
		} else if qs.Get("event") != "" {
			selectedEventLabel := qs.Get("event")
			filteredGameEvents := FilterGameEvents(cache.GameEvents, selectedEventLabel)
			messages := GenerateGameEventMessages(filteredGameEvents)

			replyMessageCall := client.ReplyMessage(event.ReplyToken, messages...)

			if _, err := replyMessageCall.Do(); err != nil {
			}
		} else if qs.Get("graphics") != "" && qs.Get("tweetId") == "" {
			selectedTwitterUser := qs.Get("graphics")
			selectedUserTweets := FindUserTweets(cache.TweetList, selectedTwitterUser)
			messages := GenerateGraphicCatalogMessages(selectedUserTweets)

			replyMessageCall := client.ReplyMessage(event.ReplyToken, messages...)

			if _, err := replyMessageCall.Do(); err != nil {
			}
		} else if qs.Get("graphics") != "" && qs.Get("tweetId") != "" {
			selectedTwitterUser := qs.Get("graphics")
			selectedTweetID := qs.Get("tweetId")
			selectedUserTweets := FindUserTweets(cache.TweetList, selectedTwitterUser)
			selectedTweet := FindTweet(selectedUserTweets.Tweets, selectedTweetID)
			messages := GenerateGraphicDetailMessages(selectedTweet, selectedTwitterUser)

			replyMessageCall := client.ReplyMessage(event.ReplyToken, messages...)

			if _, err := replyMessageCall.Do(); err != nil {
			}
		}
	}

	fmt.Fprint(w, "ok")
}
