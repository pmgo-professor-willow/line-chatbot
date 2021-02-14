package utils

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

// GetIcon
func GetIcon(selectType string) linebot.FlexComponent {
	withoutFlex := 0

	fileName := "POKEMON_TYPE_"

	switch selectType {
	case "蟲":
		fileName += "BUG"
		break
	case "惡":
		fileName += "DARK"
		break
	case "龍":
		fileName += "DRAGON"
		break
	case "電":
		fileName += "ELECTRIC"
		break
	case "妖精":
		fileName += "FAIRY"
		break
	case "火":
		fileName += "FIRE"
		break
	case "格鬥":
		fileName += "FIGHTING"
		break
	case "飛行":
		fileName += "FLYING"
		break
	case "幽靈":
		fileName += "GHOST"
		break
	case "草":
		fileName += "GRASS"
		break
	case "地面":
		fileName += "GROUND"
		break
	case "冰":
		fileName += "ICE"
		break
	case "一般":
		fileName += "NORMAL"
		break
	case "毒":
		fileName += "POISON"
		break
	case "超能力":
		fileName += "PSYCHIC"
		break
	case "岩石":
		fileName += "ROCK"
		break
	case "鋼":
		fileName += "STEEL"
		break
	case "水":
		fileName += "WATER"
		break
	}

	url := fmt.Sprintf("https://raw.githubusercontent.com/pmgo-professor-willow/line-chatbot/main/assets/types/%s.png", fileName)

	return &linebot.ImageComponent{
		Type:    linebot.FlexComponentTypeImage,
		Size:    "20px",
		URL:     url,
		Align:   linebot.FlexComponentAlignTypeCenter,
		Gravity: linebot.FlexComponentGravityTypeCenter,
		Flex:    &withoutFlex,
		Margin:  linebot.FlexComponentMarginTypeXs,
	}
}
