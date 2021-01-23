package functions

import (
	"fmt"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

// GenerateEggMessages converts eggs to LINE flex messages
func GenerateEggMessages(eggs []Egg, eggCategory string) []linebot.SendingMessage {
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
func GenerateEggBubbleMessage(eggs []Egg, eggCategory string) *linebot.BubbleContainer {
	eggRows := funk.Chunk(eggs, 3).([][]Egg)

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
				},
			},
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: funk.Map(eggRows, func(eggRow []Egg) linebot.FlexComponent {
				return &linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: funk.Map(eggRow, func(egg Egg) linebot.FlexComponent {
						return &linebot.BoxComponent{
							Type:     linebot.FlexComponentTypeBox,
							Layout:   linebot.FlexBoxLayoutTypeVertical,
							Contents: GenerateEggFlexComponent(egg),
						}
					}).([]linebot.FlexComponent),
				}
			}).([]linebot.FlexComponent),
		},
	}
}

// GenerateEggFlexComponent converts eggs to LINE bubble message
func GenerateEggFlexComponent(egg Egg) []linebot.FlexComponent {
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
		},
	}
}
