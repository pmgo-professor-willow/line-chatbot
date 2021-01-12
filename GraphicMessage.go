package functions

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

func GenerateGraphicMessage(twitterUserName string, tweet TweetData) *linebot.ImageCarouselColumn {
	return linebot.NewImageCarouselColumn(
		tweet.Media.URL,
		&linebot.PostbackAction{
			Label: "檢視完整資訊",
			Data: fmt.Sprintf(
				"graphics=%s&tweetId=%s",
				twitterUserName, tweet.Id,
			),
		},
	)
}
