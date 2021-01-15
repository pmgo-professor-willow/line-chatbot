package functions

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

// GenerateGraphicMessage converts user tweets to LINE template messages
func GenerateGraphicMessage(userTweets UserTweets) *linebot.TemplateMessage {
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
