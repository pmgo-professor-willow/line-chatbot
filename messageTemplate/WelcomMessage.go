package messageTemplate

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

// GenerateWelcomMessages sends to LINE image messages
func GenerateWelcomMessages(trainerName string) []linebot.SendingMessage {
	return []linebot.SendingMessage{
		linebot.NewTextMessage(
			fmt.Sprintf(
				"%s 初次見面，我是維羅博士。\n\n我致力於調查協助訓練家完成各項研究任務，如果遇到了任何困難都歡迎到我的實驗室詢問。",
				trainerName,
			),
		),
		linebot.NewImageMessage(
			"https://raw.githubusercontent.com/pmgo-professor-willow/line-chatbot/main/assets/tutorial-01.png",
			"https://raw.githubusercontent.com/pmgo-professor-willow/line-chatbot/main/assets/tutorial-01.png",
		),
	}
}
