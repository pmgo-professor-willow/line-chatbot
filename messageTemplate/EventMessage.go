package messageTemplate

import (
	"fmt"
	"os"
	"time"

	gd "pmgo-professor-willow/lineChatbot/gamedata"
	"pmgo-professor-willow/lineChatbot/messageTemplate/utils"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

// RemainingTime is converted by time.Duration
type RemainingTime struct {
	Days    int
	Hours   int
	Minutes int
}

// GenerateEventListMessages sends LINE quick reply messages
func GenerateEventListMessages() []linebot.SendingMessage {
	return []linebot.SendingMessage{
		linebot.NewTextMessage(
			"想知道進行中還是即將結束的活動？ (英文版資訊)",
		).WithQuickReplies(
			linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					&linebot.PostbackAction{
						Label:       "進行中的活動",
						Data:        "event=current",
						DisplayText: "請列出進行中的活動。",
					},
				),
				linebot.NewQuickReplyButton(
					"",
					&linebot.PostbackAction{
						Label:       "即將開始的活動",
						Data:        "event=upcoming",
						DisplayText: "請列出即將開始的活動。",
					},
				),
			),
		),
	}
}

// GenerateEventMessages converts game events to LINE flex messages
func GenerateEventMessages(gameEvents []gd.Event) []linebot.SendingMessage {
	if utils.IsEmpty(gameEvents) {
		return utils.GenerateEmptyReasonMessage()
	}

	eventChunks := funk.Chunk(gameEvents, 10).([][]gd.Event)

	return funk.Map(eventChunks, func(eventChunk []gd.Event) linebot.SendingMessage {
		return linebot.NewFlexMessage(
			"進行中的活動",
			&linebot.CarouselContainer{
				Type:     linebot.FlexContainerTypeCarousel,
				Contents: funk.Map(eventChunk, GenerateEventBubbleMessage).([]*linebot.BubbleContainer),
			},
		)
	}).([]linebot.SendingMessage)
}

// GenerateEventBubbleMessage converts game event to LINE bubble message
func GenerateEventBubbleMessage(event gd.Event) *linebot.BubbleContainer {
	maxFlex := 10
	withoutFlex := 0
	remainingText := "尚未公布相關時間"

	if event.Label == "upcoming" && event.StartTime != "" {
		var startTime time.Time
		if event.IsLocaleTime {
			loc, _ := time.LoadLocation(os.Getenv("TIMEZONE_LOCATION"))
			startTime, _ = time.ParseInLocation("2006-01-02T15:04:05Z", event.StartTime, loc)
		} else {
			startTime, _ = time.Parse(time.RFC3339, event.StartTime)
		}
		duration := startTime.Sub(time.Now())
		remaining := RemainingTime{
			Days:    int(duration.Hours()) / 24,
			Hours:   int(duration.Hours()) % 24,
			Minutes: int(duration.Minutes()) % 60,
		}
		remainingText = fmt.Sprintf(
			"%d 天 %d 小時 %d 分鐘後開始",
			remaining.Days, remaining.Hours, remaining.Minutes,
		)
	}

	if event.Label == "current" && event.EndTime != "" {
		var endTime time.Time
		if event.IsLocaleTime {
			loc, _ := time.LoadLocation(os.Getenv("TIMEZONE_LOCATION"))
			endTime, _ = time.ParseInLocation("2006-01-02T15:04:05Z", event.EndTime, loc)
		} else {
			endTime, _ = time.Parse(time.RFC3339, event.EndTime)
		}
		duration := endTime.Sub(time.Now())
		remaining := RemainingTime{
			Days:    int(duration.Hours()) / 24,
			Hours:   int(duration.Hours()) % 24,
			Minutes: int(duration.Minutes()) % 60,
		}
		remainingText = fmt.Sprintf(
			"剩餘 %d 天 %d 小時 %d 分鐘",
			remaining.Days, remaining.Hours, remaining.Minutes,
		)
	}

	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeMega,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeText,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type: linebot.FlexComponentTypeText,
					Text: event.Title,
				},
			},
		},
		Hero: &linebot.ImageComponent{
			Type:        linebot.FlexComponentTypeImage,
			Size:        linebot.FlexImageSizeTypeFull,
			URL:         event.ImageURL,
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
						// Empty text for padding.
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "_",
							Size:  linebot.FlexTextSizeTypeLg,
							Flex:  &maxFlex,
							Color: "#FFFFFF",
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  remainingText,
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
					Type:   linebot.FlexComponentTypeButton,
					Style:  linebot.FlexButtonStyleTypeLink,
					Action: linebot.NewURIAction("檢視活動資訊", event.Link),
				},
			},
		},
	}
}
