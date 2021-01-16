package functions

import (
	"fmt"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

// PastTime is converted by time.Duration
type PastTime struct {
	Days    int
	Hours   int
	Minutes int
}

// GenerateVideoChannelsMessages converts channels to LINE quick reply messages
func GenerateVideoChannelsMessages(channels []Channel) []linebot.SendingMessage {
	return []linebot.SendingMessage{
		linebot.NewTextMessage(
			"你想要撥放誰的影片？",
		).WithQuickReplies(
			linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					&linebot.PostbackAction{
						Label:       channels[0].Name,
						Data:        fmt.Sprintf("channel=%s", channels[0].Name),
						DisplayText: fmt.Sprintf("%s 的影片", channels[0].Name),
					},
				),
			),
		),
	}
}

// GenerateVideosMessages converts user tweets to LINE template messages
func GenerateVideosMessages(channel Channel) []linebot.SendingMessage {
	videoChunks := funk.Chunk(channel.Videos, 10).([][]Video)

	return funk.Map(videoChunks, func(videoChunk []Video) linebot.SendingMessage {
		return linebot.NewFlexMessage(
			fmt.Sprintf("%s 的 Pokemon GO 影片", channel.Name),
			&linebot.CarouselContainer{
				Type: linebot.FlexContainerTypeCarousel,
				Contents: funk.Map(
					videoChunk,
					GenerateVideoBubbleMessage,
				).([]*linebot.BubbleContainer),
			},
		)
	}).([]linebot.SendingMessage)
}

// GenerateVideoBubbleMessage converts video to LINE bubble message
func GenerateVideoBubbleMessage(video Video) *linebot.BubbleContainer {
	minFlex := 1
	withoutFlex := 0

	publishedAt, _ := time.Parse(time.RFC3339, video.PublishedAt)
	duration := time.Now().Sub(publishedAt)
	pastTime := PastTime{
		Days:    int(duration.Hours()) / 24,
		Hours:   int(duration.Hours()) % 24,
		Minutes: int(duration.Minutes()) % 60,
	}
	pastTimeText := "未知時間"
	if pastTime.Days > 0 {
		pastTimeText = fmt.Sprintf("%d 天前", pastTime.Days)
	} else if pastTime.Hours > 0 {
		pastTimeText = fmt.Sprintf("%d 小時前", pastTime.Hours)
	} else if pastTime.Minutes > 0 {
		pastTimeText = fmt.Sprintf("%d 分鐘前", pastTime.Minutes)
	}

	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeMega,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						&linebot.TextComponent{
							Type: linebot.FlexComponentTypeText,
							Text: video.Title,
						},
					},
				},
			},
		},
		Hero: &linebot.ImageComponent{
			Type:        linebot.FlexComponentTypeImage,
			Size:        linebot.FlexImageSizeTypeFull,
			URL:         video.ThumbnailUrl,
			AspectRatio: "760:572",
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
							Text:  video.ChannelTitle,
							Color: "#6C757D",
							Flex:  &minFlex,
							Align: linebot.FlexComponentAlignTypeStart,
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  pastTimeText,
							Color: "#6C757D",
							Flex:  &withoutFlex,
							Align: linebot.FlexComponentAlignTypeEnd,
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
					Type:  linebot.FlexComponentTypeButton,
					Style: linebot.FlexButtonStyleTypeLink,
					Action: &linebot.URIAction{
						Label: "播放 YouTube 影片",
						URI:   video.URL,
					},
				},
			},
		},
	}
}
