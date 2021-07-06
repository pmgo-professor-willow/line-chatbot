package functions

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	gd "pmgo-professor-willow/lineChatbot/gamedata"
	mt "pmgo-professor-willow/lineChatbot/messageTemplate"
	mtUtils "pmgo-professor-willow/lineChatbot/messageTemplate/utils"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/thoas/go-funk"
)

var botInfo *linebot.BotInfoResponse
var cache = gd.GetCache()

// WebhookFunction is base LINE webhook entry
func WebhookFunction(w http.ResponseWriter, req *http.Request) {
	// LINE messaging API client
	var client, _ = linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		linebot.WithHTTPClient(&http.Client{}),
	)

	if botInfo == nil {
		currentBotInfo, _ := client.GetBotInfo().Do()
		botInfo = currentBotInfo
	}

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
	messages := []linebot.SendingMessage{}

	if qs.Get("officialPost") != "" {
		// Refresh cache about data from cloud.
		if time.Since(cache.OfficialPostsUpdatedAt).Minutes() > 1 {
			gd.LoadOfficialPosts(&cache.OfficialPosts)
			cache.OfficialPostsUpdatedAt = time.Now()
		}

		selectedPost := qs.Get("officialPost")
		if selectedPost == "list" {
			messages = mt.GenerateOfficialPostMessages(cache.OfficialPosts)
		}

		// If empty.
		if mtUtils.IsEmpty(messages) {
			messages = mtUtils.GenerateEmptyReasonMessage()
		}
	} else if qs.Get("raidTier") != "" {
		// Refresh cache about data from cloud.
		if time.Since(cache.RaidBossesUpdatedAt).Minutes() > 1 {
			gd.LoadRaidBosses(&cache.RaidBosses)
			cache.RaidBossesUpdatedAt = time.Now()
		}

		selectedRaidTierRaw := qs.Get("raidTier")
		if selectedRaidTierRaw == "list" {
			messages = mt.GenerateRaidTierListMessages()
		} else {
			selectedRaidTiers := strings.Split(selectedRaidTierRaw, ",")

			funk.ForEach(selectedRaidTiers, func(selectedRaidTier string) {
				selectedRaidBosses := gd.FilterdRaidBosses(cache.RaidBosses, selectedRaidTier)
				messages = append(
					messages,
					mt.GenerateRaidBossMessages(selectedRaidBosses, selectedRaidTier)...,
				)
			})

			// If empty.
			if mtUtils.IsEmpty(messages) {
				messages = mtUtils.GenerateEmptyReasonMessage()
			}
		}
	} else if qs.Get("egg") != "" {
		// Refresh cache about data from cloud.
		if time.Since(cache.EggsUpdatedAt).Minutes() > 1 {
			gd.LoadEggs(&cache.Eggs)
			cache.EggsUpdatedAt = time.Now()
		}

		selectedEggCategory := qs.Get("egg")
		if selectedEggCategory == "list" {
			messages = mt.GenerateEggListMessages()
		} else {
			selectedEggs := gd.FilterdEggs(cache.Eggs, selectedEggCategory)
			messages = mt.GenerateEggMessages(selectedEggs, selectedEggCategory)
		}
	} else if qs.Get("research") != "" {
		// Refresh cache about data from cloud.
		if time.Since(cache.ResearchesUpdatedAt).Minutes() > 1 {
			gd.LoadResearches(&cache.Researches)
			cache.ResearchesUpdatedAt = time.Now()
		}

		selectedLabel := qs.Get("research")
		if selectedLabel == "list" {
			messages = mt.GenerateResearchListMessages()
		} else {
			filteredResearches := gd.FilterResearches(cache.Researches, selectedLabel)
			messages = mt.GenerateResearchMessages(filteredResearches)
		}
	} else if qs.Get("rocketInvasion") != "" {
		// Refresh cache about data from cloud.
		if time.Since(cache.RocketInvasionsUpdatedAt).Minutes() > 1 {
			gd.LoadRocketInvasions(&cache.RocketInvasions)
			cache.RocketInvasionsUpdatedAt = time.Now()
		}

		selectedLabel := qs.Get("rocketInvasion")
		if selectedLabel == "list" {
			messages = mt.GenerateRocketInvasionListMessages()
		} else if selectedLabel == "grunt" || selectedLabel == "special" {
			filteredRocketInvasions := gd.FilterRocketInvasions(cache.RocketInvasions, selectedLabel)
			messages = mt.GenerateRocketInvasionMessage(filteredRocketInvasions)
		} else {
			selectedCategory := selectedLabel
			foundRocketInvasion := gd.FindRocketInvasion(cache.RocketInvasions, selectedCategory)
			messages = mt.GenerateRocketInvasionWeaknessMessage(foundRocketInvasion)
		}
	} else if qs.Get("event") != "" {
		// Refresh cache about data from cloud.
		if time.Since(cache.EventsUpdatedAt).Minutes() > 1 {
			gd.LoadEvents(&cache.Events)
			cache.EventsUpdatedAt = time.Now()
		}

		selectedEventRaw := qs.Get("event")
		if selectedEventRaw == "list" {
			messages = mt.GenerateEventListMessages()
		} else {
			selectedEvents := strings.Split(selectedEventRaw, ",")
			selectedEventLabel := selectedEvents[0]
			var filteredEvents []gd.Event

			if selectedEventRaw == "upcoming" {
				messages = mt.GenerateEventTypeListMessages(selectedEventLabel)
			} else if len(selectedEvents) == 1 {
				filteredEvents = gd.FilterEvents(cache.Events, selectedEventLabel, nil)
				messages = mt.GenerateEventMessages(filteredEvents)
			} else {
				selectedEventType := selectedEvents[1]
				filteredEvents = gd.FilterEvents(cache.Events, selectedEventLabel, selectedEventType)
				messages = mt.GenerateEventMessages(filteredEvents)
			}
		}
	} else if qs.Get("graphics") != "" {
		// Refresh cache about data from cloud.
		if time.Since(cache.TweetListUpdatedAt).Minutes() > 1 {
			gd.LoadTweetList(&cache.TweetList)
			cache.TweetListUpdatedAt = time.Now()
		}

		selectedTwitterUser := qs.Get("graphics")
		selectedTweetID := qs.Get("tweetId")
		selectedUserTweets := gd.FindUserTweets(cache.TweetList, selectedTwitterUser)
		if selectedTweetID == "" {
			messages = mt.GenerateGraphicCatalogMessages(selectedUserTweets)
		} else {
			selectedTweet := gd.FindTweet(selectedUserTweets.Tweets, selectedTweetID)
			messages = mt.GenerateGraphicDetailMessages(selectedTweet, selectedTwitterUser)
		}
	} else if qs.Get("channel") != "" {
		// Refresh cache about data from cloud.
		if time.Since(cache.ChannelsUpdatedAt).Minutes() > 1 {
			gd.LoadChannels(&cache.Channels)
			cache.ChannelsUpdatedAt = time.Now()
		}

		selectedChannelName := qs.Get("channel")
		if selectedChannelName == "list" {
			messages = mt.GenerateVideoChannelsMessages(cache.Channels)
		} else {
			selectedChannel := gd.FindChannel(cache.Channels, selectedChannelName)
			messages = mt.GenerateVideosMessages(selectedChannel)
		}
	} else if qs.Get("faq") != "" {
		selectedQuestion := qs.Get("faq")
		if selectedQuestion == "list" {
			messages = mt.GenerateQuestionListMessages(botInfo.BasicID)
		} else {
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
