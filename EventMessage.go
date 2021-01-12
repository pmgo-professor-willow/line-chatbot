package functions

import (
	"fmt"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

type RemainingTime struct {
	Days    int
	Hours   int
	Minutes int
}

func GenerateEventBubbleMessage(event GameEvent) *linebot.BubbleContainer {
	flex := 0
	remainingText := "尚未公布結束時間"

	if event.EndTime != "" {
		endTime, _ := time.Parse(time.RFC3339, event.EndTime)
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
			URL:         event.ImageUrl,
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
							Text:  remainingText,
							Color: "#6C757D",
							Flex:  &flex,
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
