package functions

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

// GenerateGraphicCatalogMessages converts user tweets to LINE template messages
func GenerateGraphicCatalogMessages(userTweets UserTweets) []linebot.SendingMessage {
	tweetChunks := funk.Chunk(userTweets.Tweets, 10).([][]TweetData)

	return funk.Map(tweetChunks, func(tweetChunk []TweetData) linebot.SendingMessage {
		return linebot.NewTemplateMessage(
			"近期的活動圖文資訊",
			&linebot.ImageCarouselTemplate{
				Columns: funk.Map(
					userTweets.Tweets,
					func(tweet TweetData) *linebot.ImageCarouselColumn {
						return GenerateGraphicColumn(userTweets.Name, tweet)
					},
				).([]*linebot.ImageCarouselColumn),
			},
		)
	}).([]linebot.SendingMessage)
}

// GenerateGraphicDetailMessages converts tweet to LINE messages
func GenerateGraphicDetailMessages(tweet TweetData, userName string) []linebot.SendingMessage {
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
func GenerateGraphicColumn(twitterUserName string, tweet TweetData) *linebot.ImageCarouselColumn {
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
