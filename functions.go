package functions

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

func WebhookFunction(w http.ResponseWriter, r *http.Request) {
	cache := GetCache()

	// Refresh cache about data from cloud.
	if time.Since(cache.UpdatedAt).Minutes() > 1 {
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

		if qs.Get("event") != "" {
			selectedEventLabel := qs.Get("event")
			eventChunks := funk.Chunk(funk.Filter(cache.GameEvents, func(gameEvent GameEvent) bool {
				isCurrentEvent := gameEvent.Label == selectedEventLabel

				isInProgress := false
				if gameEvent.StartTime != "" && gameEvent.EndTime != "" {
					startTime, _ := time.Parse(time.RFC3339, gameEvent.StartTime)
					endTime, _ := time.Parse(time.RFC3339, gameEvent.EndTime)
					isInProgress = int(time.Now().Sub(startTime).Minutes()) > 0 && int(endTime.Sub(time.Now()).Minutes()) > 0
				} else if gameEvent.EndTime != "" {
					endTime, _ := time.Parse(time.RFC3339, gameEvent.EndTime)
					isInProgress = int(endTime.Sub(time.Now()).Minutes()) > 0
				} else if gameEvent.StartTime == "" && gameEvent.EndTime == "" {
					isInProgress = true
				}

				return isCurrentEvent && isInProgress
			}).([]GameEvent), 10).([][]GameEvent)

			eventChunkMessages := funk.Map(eventChunks, func(eventChunk []GameEvent) linebot.SendingMessage {
				return linebot.NewFlexMessage(
					"進行中的活動",
					&linebot.CarouselContainer{
						Type:     linebot.FlexContainerTypeCarousel,
						Contents: funk.Map(eventChunk, GenerateEventBubbleMessage).([]*linebot.BubbleContainer),
					},
				)
			}).([]linebot.SendingMessage)

			replyMessageCall := client.ReplyMessage(event.ReplyToken, eventChunkMessages...)

			if _, err := replyMessageCall.Do(); err != nil {
			}
		} else if qs.Get("graphics") != "" && qs.Get("tweetId") == "" {
			selectedTwitterUser := qs.Get("graphics")
			selectedUserTweets := funk.Find(cache.TweetList, func(userTweets UserTweets) bool {
				return userTweets.Name == selectedTwitterUser
			}).(UserTweets)

			message := linebot.NewTemplateMessage(
				"近期的活動圖文資訊",
				&linebot.ImageCarouselTemplate{
					Columns: funk.Map(
						selectedUserTweets.Tweets,
						func(tweet TweetData) *linebot.ImageCarouselColumn {
							return GenerateGraphicMessage(selectedUserTweets.Name, tweet)
						},
					).([]*linebot.ImageCarouselColumn),
				},
			)

			if _, err := client.ReplyMessage(event.ReplyToken, message).Do(); err != nil {
			}
		} else if qs.Get("graphics") != "" && qs.Get("tweetId") != "" {
			selectedTwitterUser := qs.Get("graphics")
			selectedTweetId := qs.Get("tweetId")
			selectedUserTweets := funk.Find(cache.TweetList, func(userTweets UserTweets) bool {
				return userTweets.Name == selectedTwitterUser
			}).(UserTweets)
			selectedTweet := funk.Find(selectedUserTweets.Tweets, func(tweet TweetData) bool {
				return tweet.Id == selectedTweetId
			}).(TweetData)

			replyMessageCall := client.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage(selectedTweet.Text),
				linebot.NewImageMessage(selectedTweet.Media.URL, selectedTweet.Media.URL),
				linebot.NewTextMessage(fmt.Sprintf(
					"以上圖文訊息由 %s 整理",
					selectedTwitterUser,
				)),
			)

			if _, err := replyMessageCall.Do(); err != nil {
			}
		}
	}

	fmt.Fprint(w, "ok")
}
