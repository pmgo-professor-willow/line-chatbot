package messageTemplate

import (
	"fmt"

	gd "pmgo-professor-willow/lineChatbot/gamedata"
	"pmgo-professor-willow/lineChatbot/messageTemplate/utils"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/thoas/go-funk"
)

// GenerateGraphicCatalogMessages converts user Instagram posts to LINE template messages
func GenerateGraphicCatalogMessages(instagramPosts gd.UserInstagramPosts) []linebot.SendingMessage {
	if utils.IsEmpty(instagramPosts) {
		return utils.GenerateEmptyReasonMessage()
	}

	postChunks := funk.Chunk(instagramPosts.Posts, 10).([][]gd.InstagramPostData)

	return funk.Map(postChunks, func(postChunk []gd.InstagramPostData) linebot.SendingMessage {
		return linebot.NewTemplateMessage(
			"近期的活動圖文資訊",
			&linebot.ImageCarouselTemplate{
				Columns: funk.Map(
					instagramPosts.Posts,
					func(post gd.InstagramPostData) *linebot.ImageCarouselColumn {
						return GenerateGraphicColumn(instagramPosts.Username, post)
					},
				).([]*linebot.ImageCarouselColumn),
			},
		)
	}).([]linebot.SendingMessage)
}

// GenerateGraphicDetailMessages converts Instagram post to LINE messages
func GenerateGraphicDetailMessages(post gd.InstagramPostData, username string) []linebot.SendingMessage {
	text := fmt.Sprintf(
		"此則圖文訊息由 %s 整理\n\n%s",
		username,
		post.Text,
	)

	messages := []linebot.SendingMessage{
		linebot.NewTextMessage(text),
	}

	for _, media := range post.MediaList {
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

// GenerateGraphicColumn converts user Instagram post to LINE column
func GenerateGraphicColumn(instagramUsername string, post gd.InstagramPostData) *linebot.ImageCarouselColumn {
	return linebot.NewImageCarouselColumn(
		post.MediaList[0].URL,
		&linebot.PostbackAction{
			Label: "檢視完整資訊",
			Data: fmt.Sprintf(
				"graphics=%s&instagramPostId=%s",
				instagramUsername, post.ID,
			),
		},
	)
}
