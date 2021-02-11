package messageTemplate

import (
	"fmt"
	"strings"

	gd "pmgo-professor-willow/lineChatbot/gamedata"
	"pmgo-professor-willow/lineChatbot/messageTemplate/utils"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

// GenerateEggMessages converts eggs to LINE flex messages
func GenerateEggMessages(eggs []gd.Egg, eggCategory string) []linebot.SendingMessage {
	if utils.IsEmpty(eggs) {
		return utils.GenerateEmptyReasonMessage()
	}

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
	appendDummyCount := (columnCount - len(eggs)%columnCount) % columnCount
	eggsWithDummies := eggs
	for i := 0; i < appendDummyCount; i++ {
		eggsWithDummies = append(eggsWithDummies, gd.CreateDummyEgg())
	}

	eggRows := funk.Chunk(eggsWithDummies, columnCount).([][]gd.Egg)

	rowContents := funk.Map(eggRows, func(eggRow []gd.Egg) linebot.FlexComponent {
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
	}).([]linebot.FlexComponent)

	// Add padding
	rowContents = append(
		rowContents,
		&linebot.SpacerComponent{
			Size: linebot.FlexSpacerSizeTypeMd,
		},
	)

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
			Type:            linebot.FlexComponentTypeBox,
			Layout:          linebot.FlexBoxLayoutTypeVertical,
			Contents:        rowContents,
			BackgroundColor: "#3D4D4D",
			Margin:          linebot.FlexComponentMarginTypeNone,
			PaddingAll:      linebot.FlexComponentPaddingTypeNone,
			PaddingStart:    linebot.FlexComponentPaddingTypeMd,
			PaddingEnd:      linebot.FlexComponentPaddingTypeMd,
		},
	}
}

// GenerateEggFlexComponent converts eggs to LINE bubble message
func GenerateEggFlexComponent(egg gd.Egg) []linebot.FlexComponent {
	maxFlex := 2
	minFlex := 1

	pokemonName := egg.Name
	// Regional pokemon.
	pokemonName = strings.Replace(pokemonName, "伽勒爾", "[伽]", 1)
	pokemonName = strings.Replace(pokemonName, "阿羅拉", "[阿]", 1)

	avatarContents := []linebot.FlexComponent{}

	// Left of avatar row
	if egg.Regional {
		avatarContents = append(
			avatarContents,
			&linebot.TextComponent{
				Type:    linebot.FlexComponentTypeText,
				Text:    "🌐",
				Size:    linebot.FlexTextSizeTypeSm,
				Align:   linebot.FlexComponentAlignTypeEnd,
				Gravity: linebot.FlexComponentGravityTypeBottom,
				Color:   "#FFFFFF",
				Flex:    &minFlex,
				Margin:  linebot.FlexComponentMarginTypeNone,
			},
		)
	} else {
		avatarContents = append(
			avatarContents,
			&linebot.FillerComponent{
				Flex: &minFlex,
			},
		)
	}

	// Center of avatar row
	avatarContents = append(
		avatarContents,
		&linebot.ImageComponent{
			Type:    linebot.FlexComponentTypeImage,
			Size:    "75px",
			URL:     egg.ImageURL,
			Align:   linebot.FlexComponentAlignTypeCenter,
			Gravity: linebot.FlexComponentGravityTypeBottom,
			Flex:    &maxFlex,
		},
	)

	// Right of avatar row
	if egg.ShinyAvailable {
		avatarContents = append(
			avatarContents,
			&linebot.TextComponent{
				Type:    linebot.FlexComponentTypeText,
				Text:    "✨",
				Size:    linebot.FlexTextSizeTypeSm,
				Align:   linebot.FlexComponentAlignTypeStart,
				Gravity: linebot.FlexComponentGravityTypeBottom,
				Color:   "#FFFFFF",
				Flex:    &minFlex,
				Margin:  linebot.FlexComponentMarginTypeNone,
			},
		)
	} else {
		avatarContents = append(
			avatarContents,
			&linebot.FillerComponent{
				Flex: &minFlex,
			},
		)
	}

	return []linebot.FlexComponent{
		&linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:     linebot.FlexComponentTypeBox,
					Layout:   linebot.FlexBoxLayoutTypeHorizontal,
					Contents: avatarContents,
				},
				&linebot.TextComponent{
					Type:  linebot.FlexComponentTypeText,
					Text:  pokemonName,
					Size:  linebot.FlexTextSizeTypeMd,
					Align: linebot.FlexComponentAlignTypeCenter,
					Color: "#FFFFFF",
				},
			},
			Margin: linebot.FlexComponentMarginTypeNone,
		},
	}
}
