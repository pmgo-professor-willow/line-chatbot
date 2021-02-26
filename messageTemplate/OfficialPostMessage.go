package messageTemplate

import (
	gd "pmgo-professor-willow/lineChatbot/gamedata"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

// GenerateOfficialPostMessages converts official posts to LINE flex messages
func GenerateOfficialPostMessages(officialPosts []gd.OfficialPost) []linebot.SendingMessage {
	postChunks := funk.Chunk(officialPosts, 10).([][]gd.OfficialPost)

	return funk.Map(postChunks, func(postChunk []gd.OfficialPost) linebot.SendingMessage {
		return linebot.NewFlexMessage(
			"近期的官方公告整理",
			&linebot.CarouselContainer{
				Type:     linebot.FlexContainerTypeCarousel,
				Contents: funk.Map(postChunk, GenerateOfficialPostBubbleMessage).([]*linebot.BubbleContainer),
			},
		)
	}).([]linebot.SendingMessage)
}

// GenerateOfficialPostBubbleMessage converts official post to LINE bubble message
func GenerateOfficialPostBubbleMessage(officialPost gd.OfficialPost) *linebot.BubbleContainer {
	minFlex := 1

	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeMega,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeText,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:  linebot.FlexComponentTypeText,
					Text:  officialPost.Date,
					Flex:  &minFlex,
					Size:  linebot.FlexTextSizeTypeLg,
					Align: linebot.FlexComponentAlignTypeCenter,
				},
			},
		},
		Hero: &linebot.ImageComponent{
			Type:        linebot.FlexComponentTypeImage,
			Size:        linebot.FlexImageSizeTypeFull,
			URL:         officialPost.CoverImageURL,
			AspectRatio: "2:1",
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  officialPost.Title,
							Color: "#6C757D",
							Flex:  &minFlex,
							Align: linebot.FlexComponentAlignTypeStart,
							Wrap:  true,
						},
					},
				},
			},
		},
		Footer: &linebot.BoxComponent{
			Type:    linebot.FlexComponentTypeBox,
			Layout:  linebot.FlexBoxLayoutTypeVertical,
			Spacing: linebot.FlexComponentSpacingTypeMd,
			Contents: []linebot.FlexComponent{
				&linebot.ButtonComponent{
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypeLink,
					Action: linebot.NewURIAction("檢視官方資訊", officialPost.Link),
				},
			},
		},
	}
}
