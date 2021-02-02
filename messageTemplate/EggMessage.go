package messageTemplate

import (
	"fmt"
	"strings"

	gd "pmgo-professor-willow/lineChatbot/gamedata"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

// GenerateEggMessages converts eggs to LINE flex messages
func GenerateEggMessages(eggs []gd.Egg, eggCategory string) []linebot.SendingMessage {
	return []linebot.SendingMessage{
		linebot.NewFlexMessage(
			fmt.Sprintf("%s 蛋可孵化出的寶可夢", eggCategory),
			&linebot.CarouselContainer{
				Type: linebot.FlexContainerTypeCarousel,
				Contents: []*linebot.BubbleContainer{
					GenerateEggBubbleMessage(eggs, eggCategory),
				},
			},
		),
	}
}

// GenerateEggBubbleMessage converts eggs to LINE bubble message
func GenerateEggBubbleMessage(eggs []gd.Egg, eggCategory string) *linebot.BubbleContainer {
	columnCount := 3
	eggsWithDummies := eggs
	for i := 0; i < columnCount-len(eggs)%columnCount; i++ {
		eggsWithDummies = append(eggsWithDummies, gd.CreateDummyEgg())
	}

	eggRows := funk.Chunk(eggsWithDummies, columnCount).([][]gd.Egg)

	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeGiga,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeText,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:  linebot.FlexComponentTypeText,
					Text:  fmt.Sprintf("%s 蛋可孵化出的寶可夢", eggCategory),
					Size:  linebot.FlexTextSizeTypeLg,
					Align: linebot.FlexComponentAlignTypeCenter,
					Color: "#FFFFFF",
				},
			},
			BackgroundColor: "#455F60",
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: funk.Map(eggRows, func(eggRow []gd.Egg) linebot.FlexComponent {
				return &linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: funk.Map(eggRow, func(egg gd.Egg) linebot.FlexComponent {
						return &linebot.BoxComponent{
							Type:     linebot.FlexComponentTypeBox,
							Layout:   linebot.FlexBoxLayoutTypeVertical,
							Contents: GenerateEggFlexComponent(egg),
						}
					}).([]linebot.FlexComponent),
				}
			}).([]linebot.FlexComponent),
			BackgroundColor: "#3D4D4D",
			Margin:          linebot.FlexComponentMarginTypeNone,
		},
	}
}

// GenerateEggFlexComponent converts eggs to LINE bubble message
func GenerateEggFlexComponent(egg gd.Egg) []linebot.FlexComponent {
	pokemonName := egg.Name
	// Regional pokemon.
	pokemonName = strings.Replace(pokemonName, "伽勒爾", "[伽]", 1)
	pokemonName = strings.Replace(pokemonName, "阿羅拉", "[阿]", 1)

	if egg.ShinyAvailable {
		pokemonName += " ✨"
	}

	return []linebot.FlexComponent{
		&linebot.ImageComponent{
			Type:  linebot.FlexComponentTypeImage,
			Size:  "75px",
			URL:   egg.ImageURL,
			Align: linebot.FlexComponentAlignTypeCenter,
		},
		&linebot.TextComponent{
			Type:  linebot.FlexComponentTypeText,
			Text:  pokemonName,
			Size:  linebot.FlexTextSizeTypeMd,
			Align: linebot.FlexComponentAlignTypeCenter,
			Color: "#FFFFFF",
		},
	}
}
