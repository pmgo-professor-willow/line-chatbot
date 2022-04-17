package messageTemplate

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// GenerateAnsweringMessage sends LINE messages
func GenerateAnsweringMessage(trainerName string) []linebot.SendingMessage {
	return []linebot.SendingMessage{
		linebot.NewTextMessage(
			fmt.Sprintf(
				"謝謝 %s 的建議，博士已經記錄下來了。",
				trainerName,
			),
		),
	}
}
