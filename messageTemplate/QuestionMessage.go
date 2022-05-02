package messageTemplate

import (
	"fmt"
	"net/url"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// GenerateQuestionListMessages sends LINE flex messages
func GenerateQuestionListMessages(botBasicID string) []linebot.SendingMessage {
	return []linebot.SendingMessage{
		linebot.NewTemplateMessage(
			"常見問題",
			linebot.NewCarouselTemplate(
				&linebot.CarouselColumn{
					ThumbnailImageURL: "https://raw.githubusercontent.com/pmgo-professor-willow/line-chatbot/main/assets/faq-donate.png",
					Title:             "贊助",
					Text:              "若您使用滿意，可以考慮鼓勵開發者",
					Actions: []linebot.TemplateAction{
						&linebot.URIAction{
							Label: "需要贊助的理由",
							URI:   "https://liff.line.me/1645278921-kWRPP32q/611mscwy/text/560773408578064?accountId=611mscwy",
						},
						&linebot.URIAction{
							Label: "使用台新 Richart 轉帳",
							URI:   "https://richart.tw/TSDIB_RichartWeb/RC04/RC040300?token=X6Y36lCy06A%3D",
						},
					},
				},
				&linebot.CarouselColumn{
					ThumbnailImageURL: "https://raw.githubusercontent.com/pmgo-professor-willow/line-chatbot/main/assets/faq-data.png",
					Title:             "資料相關",
					Text:              "關於團體戰、蛋池與活動等資訊",
					Actions: []linebot.TemplateAction{
						&linebot.PostbackAction{
							Label:       "資料來源與更新週期",
							Data:        "faq=dataSource",
							DisplayText: "我想知道資料來源與更新週期是？",
						},
						&linebot.PostbackAction{
							Label:       "資料正確性",
							Data:        "faq=dataAccuracy",
							DisplayText: "我想知道資料的正確性有多高？",
						},
					},
				},
				&linebot.CarouselColumn{
					ThumbnailImageURL: "https://raw.githubusercontent.com/pmgo-professor-willow/line-chatbot/main/assets/faq-misc.png",
					Title:             "其它問題",
					Text:              "關於維羅博士的運作方式與系統反饋",
					Actions: []linebot.TemplateAction{
						&linebot.PostbackAction{
							Label:       "服務完全免費",
							Data:        "faq=pricing",
							DisplayText: "我想知道這項服務是免費還是付費的？",
						},
						&linebot.PostbackAction{
							Label:       "提供建議或反饋",
							Data:        "faq=contact",
							DisplayText: "我應該如何提供對系統的建議或反饋？",
						},
					},
				},
				&linebot.CarouselColumn{
					ThumbnailImageURL: "https://raw.githubusercontent.com/pmgo-professor-willow/line-chatbot/main/assets/faq-share.png",
					Title:             "分享推廣",
					Text:              "將維羅博士介紹給更多的訓練家",
					Actions: []linebot.TemplateAction{
						&linebot.URIAction{
							Label: "將博士介紹給朋友",
							URI: fmt.Sprintf(
								"https://line.me/R/nv/recommendOA/%s",
								botBasicID,
							),
						},
						&linebot.URIAction{
							Label: "巴哈姆特討論串",
							URI:   "https://forum.gamer.com.tw/C.php?bsn=29659&snA=40930",
						},
					},
				},
			),
		),
	}
}

// GenerateQuestionMessages sends LINE flex messages
func GenerateQuestionMessages(selectedQuestion, botBasicID string) []linebot.SendingMessage {
	messages := []linebot.SendingMessage{}

	if selectedQuestion == "dataSource" {
		messages = []linebot.SendingMessage{
			linebot.NewTextMessage(
				"資料來源主以國外 Leek Duck 與 The Sliph Road 網站所彙整，維羅博士透過自動化程式進行收集。\n\n因此更新時間將以上述網站為主，而雙方資訊差異不會超過三十分鐘。",
			),
			linebot.NewTextMessage(
				"維羅博士所使用之圖片、寶可夢資訊之版權屬於 Niantic, Inc. 與 Nintendo 擁有。（部分為二創將不在此列）",
			),
		}
	} else if selectedQuestion == "dataAccuracy" {
		messages = []linebot.SendingMessage{
			linebot.NewTextMessage(
				"資料取自富有規模的國外資料站，儘管可信度相當高，若與實際遊戲內容存在差異，維羅博士不另行告知。",
			),
			linebot.NewTextMessage(
				"因地方時區因素，可能存在活動交替導致資訊落差，請各位訓練家注意。\n\n而時間倒數資訊將以台灣時區為主 (GMT+8)。",
			),
		}
	} else if selectedQuestion == "pricing" {
		messages = []linebot.SendingMessage{
			linebot.NewTextMessage(
				"維羅博士提供的功能皆為「免費」，且不會有任何廣告訊息。\n\n在使用過程中，傳輸圖片所產生的流量，請訓練家們自行注意哦！",
			),
		}
	} else if selectedQuestion == "contact" {
		messages = []linebot.SendingMessage{
			linebot.NewFlexMessage(
				"與博士聯繫",
				&linebot.BubbleContainer{
					Type: linebot.FlexContainerTypeBubble,
					Size: linebot.FlexBubbleSizeTypeMega,
					Hero: &linebot.ImageComponent{
						Type:        linebot.FlexComponentTypeImage,
						Size:        linebot.FlexImageSizeTypeFull,
						URL:         "https://raw.githubusercontent.com/pmgo-professor-willow/line-chatbot/main/assets/author.png",
						AspectRatio: "648:355",
					},
					Body: &linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeVertical,
						Contents: []linebot.FlexComponent{
							&linebot.TextComponent{
								Type:   linebot.FlexComponentTypeText,
								Text:   "如果是想要匿名留給維羅博士，請直接發送首三字為「給博士」的文字訊息。\n\n如果有任何建議都歡迎來信至博士的 Email，標題請與「維羅博士」字眼相關。",
								Color:  "#6C757D",
								Align:  linebot.FlexComponentAlignTypeStart,
								Wrap:   true,
								Margin: linebot.FlexComponentMarginTypeSm,
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
									Label: "傳送敲敲話給博士",
									URI: fmt.Sprintf(
										"https://line.me/R/oaMessage/%s/?%s",
										botBasicID,
										url.QueryEscape("給博士，"),
									),
								},
							},
							&linebot.ButtonComponent{
								Type:  linebot.FlexComponentTypeButton,
								Style: linebot.FlexButtonStyleTypeLink,
								Action: &linebot.URIAction{
									Label: "寫信給博士",
									URI: fmt.Sprintf(
										"mailto:salmon.zh.tw@gmail.com?subject=%s&body=%s",
										url.QueryEscape("訓練家給維羅博士的一封信"),
										url.QueryEscape("博士您好，"),
									),
								},
							},
						},
					},
				},
			),
		}
	}

	return messages
}
