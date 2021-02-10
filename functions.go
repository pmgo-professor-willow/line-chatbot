package functions

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	gd "pmgo-professor-willow/lineChatbot/gamedata"
	mt "pmgo-professor-willow/lineChatbot/messageTemplate"

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
		case linebot.EventTypeFollow:
			profile, _ := client.GetProfile(event.Source.UserID).Do()
			messages := mt.GenerateWelcomMessages(profile.DisplayName)

			replyMessageCall := client.ReplyMessage(event.ReplyToken, messages...)

			if _, err := replyMessageCall.Do(); err != nil {
				fmt.Println(err)
			}
			break

		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				// command: send postback message directly.
				postbackData := regexp.MustCompile(`^/(.+)$`).FindStringSubmatch(message.Text)
				if len(postbackData) != 0 {
					qs, _ := url.ParseQuery(postbackData[1])
					PostbackReply(client, event.ReplyToken, qs)
					break
				}

				// leave meassage: send message to manager.
				isLeaveMessage := regexp.MustCompile(`^給博士.+`).MatchString(message.Text)
				if isLeaveMessage {
					profile, _ := client.GetProfile(event.Source.UserID).Do()
					PushMessageToManager(
						client,
						fmt.Sprintf("%s: %s", profile.DisplayName, message.Text),
					)
					break
				}
				break
			}

		case linebot.EventTypePostback:
			qs, _ := url.ParseQuery(event.Postback.Data)
			PostbackReply(client, event.ReplyToken, qs)
			break
		}
	}

	fmt.Fprint(w, "ok")
}

// PostbackReply will reply messages by postback
func PostbackReply(client *linebot.Client, replyToken string, qs url.Values) {
	cache := gd.GetCache()
	messages := []linebot.SendingMessage{}

	// Refresh cache about data from cloud.
	if time.Since(cache.UpdatedAt).Minutes() > 1 {
		cache.RaidBosses = gd.LoadRaidBosses()
		cache.Eggs = gd.LoadEggs()
		cache.Researches = gd.LoadResearches()
		cache.RocketInvasions = gd.LoadRocketInvasions()
		cache.Events = gd.LoadEvents()
		cache.TweetList = gd.LoadTweetList()
		cache.Channels = gd.LoadChannels()
		cache.UpdatedAt = time.Now()
	}

	if qs.Get("raidTier") != "" {
		selectedRaidTier := qs.Get("raidTier")
		selectedRaidBosses := gd.FilterdRaidBosses(cache.RaidBosses, selectedRaidTier)
		messages = mt.GenerateRaidBossMessages(selectedRaidBosses, selectedRaidTier)
	} else if qs.Get("egg") != "" {
		selectedEggCategory := qs.Get("egg")
		selectedEggs := gd.FilterdEggs(cache.Eggs, selectedEggCategory)
		messages = mt.GenerateEggMessages(selectedEggs, selectedEggCategory)
	} else if qs.Get("researches") != "" {
		messages = mt.GenerateResearchMessages(cache.Researches)
	} else if qs.Get("rocketInvasion") != "" {
		selectedLabel := qs.Get("rocketInvasion")
		if selectedLabel == "list" {
			messages = mt.GenerateRocketInvasionListMessages()
		} else {
			filteredRocketInvasions := gd.FilterdRocketInvasions(cache.RocketInvasions, selectedLabel)
			messages = mt.GenerateRocketInvasionMessage(filteredRocketInvasions)
		}
	} else if qs.Get("event") != "" {
		selectedEventLabel := qs.Get("event")
		if selectedEventLabel == "list" {
			messages = mt.GenerateEventListMessages()
		} else {
			filteredEvents := gd.FilterEvents(cache.Events, selectedEventLabel)
			messages = mt.GenerateEventMessages(filteredEvents)
		}
	} else if qs.Get("graphics") != "" && qs.Get("tweetId") == "" {
		selectedTwitterUser := qs.Get("graphics")
		selectedUserTweets := gd.FindUserTweets(cache.TweetList, selectedTwitterUser)
		messages = mt.GenerateGraphicCatalogMessages(selectedUserTweets)
	} else if qs.Get("graphics") != "" && qs.Get("tweetId") != "" {
		selectedTwitterUser := qs.Get("graphics")
		selectedTweetID := qs.Get("tweetId")
		selectedUserTweets := gd.FindUserTweets(cache.TweetList, selectedTwitterUser)
		selectedTweet := gd.FindTweet(selectedUserTweets.Tweets, selectedTweetID)
		messages = mt.GenerateGraphicDetailMessages(selectedTweet, selectedTwitterUser)
	} else if qs.Get("channels") != "" {
		messages = mt.GenerateVideoChannelsMessages(cache.Channels)
	} else if qs.Get("channel") != "" {
		selectedChannelName := qs.Get("channel")
		selectedChannel := gd.FindChannel(cache.Channels, selectedChannelName)
		messages = mt.GenerateVideosMessages(selectedChannel)
	} else if qs.Get("faq") != "" {
		selectedQuestion := qs.Get("faq")
		if selectedQuestion == "list" {
			messages = mt.GenerateQuestionListMessages()
		} else {
			botInfo, _ := client.GetBotInfo().Do()
			messages = mt.GenerateQuestionMessages(selectedQuestion, botInfo.BasicID)
		}
	}

	if len(messages) > 0 {
		replyMessageCall := client.ReplyMessage(replyToken, messages...)

		if _, err := replyMessageCall.Do(); err != nil {
			fmt.Println(err)
		}
	}
}

// PushMessageToManager will push text message to manager
func PushMessageToManager(client *linebot.Client, messageText string) {
	replyMessageCall := client.PushMessage(
		os.Getenv("LINE_MANAGER_USER_ID"),
		linebot.NewTextMessage(messageText),
	)

	if _, err := replyMessageCall.Do(); err != nil {
		fmt.Println(err)
	}
}
