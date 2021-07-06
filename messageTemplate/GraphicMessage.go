package messageTemplate

import (
	"fmt"

	gd "pmgo-professor-willow/lineChatbot/gamedata"
	"pmgo-professor-willow/lineChatbot/messageTemplate/utils"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/thoas/go-funk"
)

// GenerateGraphicCatalogMessages converts user tweets to LINE template messages
func GenerateGraphicCatalogMessages(userTweets gd.UserTweets) []linebot.SendingMessage {
	if utils.IsEmpty(userTweets) {
		return utils.GenerateEmptyReasonMessage()
	}

	tweetChunks := funk.Chunk(userTweets.Tweets, 10).([][]gd.TweetData)

	return funk.Map(tweetChunks, func(tweetChunk []gd.TweetData) linebot.SendingMessage {
		return linebot.NewTemplateMessage(
			"近期的活動圖文資訊",
			&linebot.ImageCarouselTemplate{
				Columns: funk.Map(
					userTweets.Tweets,
					func(tweet gd.TweetData) *linebot.ImageCarouselColumn {
						return GenerateGraphicColumn(userTweets.Name, tweet)
					},
				).([]*linebot.ImageCarouselColumn),
			},
		)
	}).([]linebot.SendingMessage)
}

// GenerateGraphicDetailMessages converts tweet to LINE messages
func GenerateGraphicDetailMessages(tweet gd.TweetData, userName string) []linebot.SendingMessage {
	text := fmt.Sprintf(
		"此則圖文訊息由 %s 整理\n\n%s",
		userName,
		tweet.Text,
	)

	messages := []linebot.SendingMessage{
		linebot.NewTextMessage(text),
	}

	for _, media := range tweet.MediaList {
		messages = append(
			messages,
			linebot.NewImageMessage(media.URL, media.URL),
		)
	}

	// Slice messages becuase message amount limit
	if len(messages) > 5 {
		messages = messages[0:5]
	}

	return messages
}

// GenerateGraphicColumn converts user tweets to LINE column
func GenerateGraphicColumn(twitterUserName string, tweet gd.TweetData) *linebot.ImageCarouselColumn {
	return linebot.NewImageCarouselColumn(
		tweet.MediaList[0].URL,
		&linebot.PostbackAction{
			Label: "檢視完整資訊",
			Data: fmt.Sprintf(
				"graphics=%s&tweetId=%s",
				twitterUserName, tweet.ID,
			),
		},
	)
}
