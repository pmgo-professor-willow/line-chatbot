package messageTemplate

import (
	"fmt"

	gd "pmgo-professor-willow/lineChatbot/gamedata"
	"pmgo-professor-willow/lineChatbot/messageTemplate/utils"

	"github.com/line/line-bot-sdk-go/linebot"
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
	return []linebot.SendingMessage{
		linebot.NewTextMessage(tweet.Text),
		linebot.NewImageMessage(tweet.Media.URL, tweet.Media.URL),
		linebot.NewTextMessage(fmt.Sprintf(
			"以上圖文訊息由 %s 整理",
			userName,
		)),
	}
}

// GenerateGraphicColumn converts user tweets to LINE column
func GenerateGraphicColumn(twitterUserName string, tweet gd.TweetData) *linebot.ImageCarouselColumn {
	return linebot.NewImageCarouselColumn(
		tweet.Media.URL,
		&linebot.PostbackAction{
			Label: "檢視完整資訊",
			Data: fmt.Sprintf(
				"graphics=%s&tweetId=%s",
				twitterUserName, tweet.ID,
			),
		},
	)
}
