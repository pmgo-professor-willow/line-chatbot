package utils

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

// IsEmpty
func IsEmpty(sourceList interface{}) bool {
	return funk.IsEmpty(sourceList)
}

func GenerateEmptyReasonMessage() []linebot.SendingMessage {
	return []linebot.SendingMessage{
		linebot.NewTextMessage(
			"資料目前無法正確載入，或是當前無相關資料，請稍後再嘗試看看。",
		),
	}
}
