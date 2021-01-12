package functions

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

func WebhookFunction(w http.ResponseWriter, r *http.Request) {
	cache := GetCache()

	// Refresh cache about data from cloud.
	if time.Since(cache.UpdatedAt).Minutes() > 1 {
		cache.GameEvents = LoadGameEvents()
		cache.UpdatedAt = time.Now()
	}

	// LINE messaging API client
	var client, _ = linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
		linebot.WithHTTPClient(&http.Client{}),
	)

	events, _ := client.ParseRequest(r)
	for _, event := range events {
		if event.Type != linebot.EventTypePostback {
			break
		}

		qs, _ := url.ParseQuery(event.Postback.Data)

		if qs.Get("event") != "" {
			selectedEventLabel := qs.Get("event")
			eventChunks := funk.Chunk(funk.Filter(cache.GameEvents, func(gameEvent GameEvent) bool {
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
						Contents: funk.Map(eventChunk, GenerateEventBubbleMessage).([]*linebot.BubbleContainer),
					},
				)
			}).([]linebot.SendingMessage)

			if _, err := client.ReplyMessage(event.ReplyToken, eventChunkMessages...).Do(); err != nil {
			}
		}
	}

	fmt.Fprint(w, "ok")
}
