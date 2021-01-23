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
func WebhookFunction(w http.ResponseWriter, req *http.Request) {
	// LINE messaging API client
	var client, _ = linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		linebot.WithHTTPClient(&http.Client{}),
	)

	events, _ := client.ParseRequest(req)
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				qs, _ := url.ParseQuery(message.Text)
				PostbackReply(client, event.ReplyToken, qs)
			}

		case linebot.EventTypePostback:
			qs, _ := url.ParseQuery(event.Postback.Data)
			PostbackReply(client, event.ReplyToken, qs)
		}
	}

	fmt.Fprint(w, "ok")
}

// PostbackReply will reply messages by postback
func PostbackReply(client *linebot.Client, replyToken string, qs url.Values) {
	cache := GetCache()

	// Refresh cache about data from cloud.
	if time.Since(cache.UpdatedAt).Minutes() > 1 {
		cache.RaidBosses = LoadRaidBosses()
		cache.Eggs = LoadEggs()
		cache.GameEvents = LoadGameEvents()
		cache.TweetList = LoadTweetList()
		cache.Channels = LoadChannels()
		cache.UpdatedAt = time.Now()
	}

	if qs.Get("raidTier") != "" {
		selectedRaidTier := qs.Get("raidTier")
		selectedRaidBosses := FilterdRaidBosses(cache.RaidBosses, selectedRaidTier)
		messages := GenerateRaidBossMessages(selectedRaidBosses, selectedRaidTier)

		replyMessageCall := client.ReplyMessage(replyToken, messages...)

		if _, err := replyMessageCall.Do(); err != nil {
		}
	} else if qs.Get("egg") != "" {
		selectedEggCategory := qs.Get("egg")
		selectedEggs := FilterdEggs(cache.Eggs, selectedEggCategory)
		messages := GenerateEggMessages(selectedEggs, selectedEggCategory)

		replyMessageCall := client.ReplyMessage(replyToken, messages...)

		if _, err := replyMessageCall.Do(); err != nil {
		}
	} else if qs.Get("event") != "" {
		selectedEventLabel := qs.Get("event")
		filteredGameEvents := FilterGameEvents(cache.GameEvents, selectedEventLabel)
		messages := GenerateGameEventMessages(filteredGameEvents)

		replyMessageCall := client.ReplyMessage(replyToken, messages...)

		if _, err := replyMessageCall.Do(); err != nil {
		}
	} else if qs.Get("graphics") != "" && qs.Get("tweetId") == "" {
		selectedTwitterUser := qs.Get("graphics")
		selectedUserTweets := FindUserTweets(cache.TweetList, selectedTwitterUser)
		messages := GenerateGraphicCatalogMessages(selectedUserTweets)

		replyMessageCall := client.ReplyMessage(replyToken, messages...)

		if _, err := replyMessageCall.Do(); err != nil {
		}
	} else if qs.Get("graphics") != "" && qs.Get("tweetId") != "" {
		selectedTwitterUser := qs.Get("graphics")
		selectedTweetID := qs.Get("tweetId")
		selectedUserTweets := FindUserTweets(cache.TweetList, selectedTwitterUser)
		selectedTweet := FindTweet(selectedUserTweets.Tweets, selectedTweetID)
		messages := GenerateGraphicDetailMessages(selectedTweet, selectedTwitterUser)

		replyMessageCall := client.ReplyMessage(replyToken, messages...)

		if _, err := replyMessageCall.Do(); err != nil {
		}
	} else if qs.Get("channels") != "" {
		messages := GenerateVideoChannelsMessages(cache.Channels)

		replyMessageCall := client.ReplyMessage(replyToken, messages...)

		if _, err := replyMessageCall.Do(); err != nil {
		}
	} else if qs.Get("channel") != "" {
		selectedChannelName := qs.Get("channel")
		selectedChannel := FindChannel(cache.Channels, selectedChannelName)
		messages := GenerateVideosMessages(selectedChannel)

		replyMessageCall := client.ReplyMessage(replyToken, messages...)

		if _, err := replyMessageCall.Do(); err != nil {
		}
	}
}
