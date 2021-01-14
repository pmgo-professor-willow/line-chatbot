package functions

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

// WebhookFunction is base LINE webhook entry
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
			filteredGameEvents := FilterGameEvents(cache.GameEvents, selectedEventLabel)
			eventChunkMessages := GenerateGameEventMessages(filteredGameEvents)

			if _, err := client.ReplyMessage(event.ReplyToken, eventChunkMessages...).Do(); err != nil {
			}
		}
	}

	fmt.Fprint(w, "ok")
}
