package lineChatbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

type GameEvent struct {
	Title        string `json:"title"`
	Link         string `json:"link"`
	Type         string `json:"type"`
	ImageUrl     string `json:"imageUrl"`
	Label        string `json:"label"`
	IsLocaleTime bool   `json:"isLocaleTime"`
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
}

type RemainingTime struct {
	Days    int
	Hours   int
	Minutes int
}

func LoadGameEvents() []GameEvent {
	resp, fetchErr := http.Get("https://pmgo-professor-willow.github.io/data-leekduck/events.json")

	if fetchErr == nil {
		defer resp.Body.Close()
		bodyBuf, readErr := ioutil.ReadAll(resp.Body)

		if readErr == nil {
			events := []GameEvent{}
			json.Unmarshal(bodyBuf, &events)

			return events
		}
	}

	return []GameEvent{}
}

func generateEventBubbleMessage(event GameEvent) *linebot.BubbleContainer {
	var flex int = 0
	var remainingText string = "尚未公布結束時間"

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

func WebhookFunction(w http.ResponseWriter, req *http.Request) {
	lineChannelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	lineChannelAccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	gameEvents := LoadGameEvents()

	clientOption := linebot.WithHTTPClient(&http.Client{})
	bot, _ := linebot.New(lineChannelSecret, lineChannelAccessToken, clientOption)

	events, _ := bot.ParseRequest(req)
	for _, event := range events {
		if event.Type != linebot.EventTypePostback {
			break
		}

		qs, _ := url.ParseQuery(event.Postback.Data)

		if qs.Get("event") != "" {
			selectedEventLabel := qs.Get("event")
			eventChunks := funk.Chunk(funk.Filter(gameEvents, func(gameEvent GameEvent) bool {
				isCurrentEvent := gameEvent.Label == selectedEventLabel

				isInProgress := false
				if gameEvent.StartTime != "" && gameEvent.EndTime != "" {
					startTime, _ := time.Parse(time.RFC3339, gameEvent.StartTime)
					endTime, _ := time.Parse(time.RFC3339, gameEvent.EndTime)
					isInProgress = int(time.Now().Sub(startTime).Minutes()) > 0 && int(endTime.Sub(time.Now()).Minutes()) > 0
				} else if gameEvent.EndTime != "" {
					endTime, _ := time.Parse(time.RFC3339, gameEvent.EndTime)
					isInProgress = int(endTime.Sub(time.Now()).Minutes()) > 0
				} else if gameEvent.StartTime == "" && gameEvent.EndTime == "" {
					isInProgress = true
				}

				return isCurrentEvent && isInProgress
			}).([]GameEvent), 10).([][]GameEvent)

			eventChunkMessages := funk.Map(eventChunks, func(eventChunk []GameEvent) linebot.SendingMessage {
				return linebot.NewFlexMessage(
					"進行中的活動",
					&linebot.CarouselContainer{
						Type:     linebot.FlexContainerTypeCarousel,
						Contents: funk.Map(eventChunk, generateEventBubbleMessage).([]*linebot.BubbleContainer),
					},
				)
			}).([]linebot.SendingMessage)

			if _, err := bot.ReplyMessage(event.ReplyToken, eventChunkMessages...).Do(); err != nil {
			}
		}
	}

	fmt.Fprint(w, "ok")
}
