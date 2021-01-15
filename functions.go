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
			selectedRaidBosses := funk.Filter(cache.RaidBosses, func(raidBoss RaidBoss) bool {
				return raidBoss.Tier == selectedRaidTier
			})

			replyMessageCall := client.ReplyMessage(
				event.ReplyToken,
				linebot.NewFlexMessage(
					fmt.Sprintf("%s 星團體戰列表", selectedRaidTier),
					&linebot.CarouselContainer{
						Type: linebot.FlexContainerTypeCarousel,
						Contents: funk.Map(
							selectedRaidBosses,
							generateRaidBossMessage,
						).([]*linebot.BubbleContainer),
					},
				),
			)

			if _, err := replyMessageCall.Do(); err != nil {
			}

			// await client.replyMessage(replyToken, [
			//   {
			// 	type: 'flex',
			// 	altText: `${qs.raidTier} 星團體戰列表`,
			// 	contents: {
			// 	  type: 'carousel',
			// 	  contents: raidBosses.map(generateRaidBossBubbleMessage),
			// 	},
			//   },
			// ]);
		} else if qs.Get("event") != "" {
			selectedEventLabel := qs.Get("event")
			filteredGameEvents := FilterGameEvents(cache.GameEvents, selectedEventLabel)
			eventChunkMessages := GenerateGameEventMessages(filteredGameEvents)

			replyMessageCall := client.ReplyMessage(event.ReplyToken, eventChunkMessages...)

			if _, err := replyMessageCall.Do(); err != nil {
			}
		} else if qs.Get("graphics") != "" && qs.Get("tweetId") == "" {
			selectedTwitterUser := qs.Get("graphics")
			selectedUserTweets := funk.Find(cache.TweetList, func(userTweets UserTweets) bool {
				return userTweets.Name == selectedTwitterUser
			}).(UserTweets)

			message := GenerateGraphicMessage(selectedUserTweets)

			replyMessageCall := client.ReplyMessage(event.ReplyToken, message)

			if _, err := replyMessageCall.Do(); err != nil {
			}
		} else if qs.Get("graphics") != "" && qs.Get("tweetId") != "" {
			selectedTwitterUser := qs.Get("graphics")
			selectedTweetID := qs.Get("tweetId")
			selectedUserTweets := funk.Find(cache.TweetList, func(userTweets UserTweets) bool {
				return userTweets.Name == selectedTwitterUser
			}).(UserTweets)
			selectedTweet := funk.Find(selectedUserTweets.Tweets, func(tweet TweetData) bool {
				return tweet.ID == selectedTweetID
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
