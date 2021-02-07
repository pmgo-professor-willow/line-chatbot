package messageTemplate

import (
	"fmt"

	gd "pmgo-professor-willow/lineChatbot/gamedata"
	"pmgo-professor-willow/lineChatbot/messageTemplate/utils"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

// GenerateRocketInvasionListMessages sends LINE quick reply messages
func GenerateRocketInvasionListMessages() []linebot.SendingMessage {
	return []linebot.SendingMessage{
		linebot.NewTextMessage(
			"我想了解目前的火箭隊陣容。",
		).WithQuickReplies(
			linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					&linebot.PostbackAction{
						Label:       "普通手下的陣容",
						Data:        "rocketInvasion=grunt",
						DisplayText: "請列出火箭隊普通手下的陣容。",
					},
				),
				linebot.NewQuickReplyButton(
					"",
					&linebot.PostbackAction{
						Label:       "幹部與特殊角色的陣容",
						Data:        "rocketInvasion=special",
						DisplayText: "請列出火箭隊幹部與特殊角色的陣容。",
					},
				),
			),
		),
	}
}

// GenerateRocketInvasionMessage converts rocket invasions to LINE flex messages
func GenerateRocketInvasionMessage(rocketInvasions []gd.RocketInvasion) []linebot.SendingMessage {
	if utils.IsEmpty(rocketInvasions) {
		return utils.GenerateEmptyReasonMessage()
	}

	rocketInvasionChunks := funk.Chunk(rocketInvasions, 5).([][]gd.RocketInvasion)

	return funk.Map(rocketInvasionChunks, func(rocketInvasionChunk []gd.RocketInvasion) linebot.SendingMessage {
		return linebot.NewFlexMessage(
			fmt.Sprintf("火箭隊成員與幹部陣容"),
			&linebot.CarouselContainer{
				Type: linebot.FlexContainerTypeCarousel,
				Contents: funk.Map(
					rocketInvasionChunk,
					GenerateRocketInvasionBubbleMessage,
				).([]*linebot.BubbleContainer),
			},
		)
	}).([]linebot.SendingMessage)
}

// GenerateRocketInvasionBubbleMessage converts rocket invasions to LINE bubble message
func GenerateRocketInvasionBubbleMessage(rocketInvasion gd.RocketInvasion) *linebot.BubbleContainer {
	maxFlex := 3
	minFlex := 1

	titleText := ""

	if rocketInvasion.IsSpecial {
		titleText = rocketInvasion.Category
	} else {
		titleText = fmt.Sprintf("手下 (%s)", rocketInvasion.Category)
	}

	titleContent := &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeVertical,
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:   linebot.FlexComponentTypeText,
				Text:   titleText,
				Size:   linebot.FlexTextSizeTypeLg,
				Align:  linebot.FlexComponentAlignTypeStart,
				Margin: linebot.FlexComponentMarginTypeNone,
				Color:  "#FFFFFF",
				Weight: linebot.FlexTextWeightTypeBold,
			},
		},
		Flex: &maxFlex,
	}

	// Append quote if exist
	if rocketInvasion.Quote != "" {
		titleContent.Contents = append(
			titleContent.Contents,
			&linebot.TextComponent{
				Type:   linebot.FlexComponentTypeText,
				Text:   rocketInvasion.Quote,
				Size:   linebot.FlexTextSizeTypeMd,
				Align:  linebot.FlexComponentAlignTypeStart,
				Margin: linebot.FlexComponentMarginTypeNone,
				Color:  "#CDCDCD",
				Weight: linebot.FlexTextWeightTypeBold,
				Wrap:   true,
			},
		)
	}

	headerContents := []linebot.FlexComponent{
		&linebot.ImageComponent{
			Type:    linebot.FlexComponentTypeImage,
			Size:    "100px",
			URL:     rocketInvasion.CharacterImageURL,
			Align:   linebot.FlexComponentAlignTypeCenter,
			Gravity: linebot.FlexComponentGravityTypeBottom,
			Margin:  linebot.FlexComponentMarginTypeNone,
			Flex:    &minFlex,
		},
		titleContent,
	}

	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeGiga,
		Header: &linebot.BoxComponent{
			Type:            linebot.FlexComponentTypeBox,
			Layout:          linebot.FlexBoxLayoutTypeHorizontal,
			Contents:        headerContents,
			BackgroundColor: "#455F60",
			PaddingBottom:   linebot.FlexComponentPaddingTypeNone,
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:     linebot.FlexComponentTypeBox,
							Layout:   linebot.FlexBoxLayoutTypeVertical,
							Contents: GenerateLineupPokemonsFlexComponent(rocketInvasion.LineupPokemons, 1),
							Margin:   linebot.FlexComponentMarginTypeXs,
							Flex:     &minFlex,
						},
						&linebot.SeparatorComponent{
							Type:   linebot.FlexComponentTypeSeparator,
							Color:  "#CDCDCD",
							Margin: linebot.FlexComponentMarginTypeNone,
						},
						&linebot.BoxComponent{
							Type:     linebot.FlexComponentTypeBox,
							Layout:   linebot.FlexBoxLayoutTypeVertical,
							Contents: GenerateLineupPokemonsFlexComponent(rocketInvasion.LineupPokemons, 2),
							Margin:   linebot.FlexComponentMarginTypeXs,
							Flex:     &minFlex,
						},
						&linebot.SeparatorComponent{
							Type:   linebot.FlexComponentTypeSeparator,
							Color:  "#CDCDCD",
							Margin: linebot.FlexComponentMarginTypeNone,
						},
						&linebot.BoxComponent{
							Type:     linebot.FlexComponentTypeBox,
							Layout:   linebot.FlexBoxLayoutTypeVertical,
							Contents: GenerateLineupPokemonsFlexComponent(rocketInvasion.LineupPokemons, 3),
							Margin:   linebot.FlexComponentMarginTypeXs,
							Flex:     &minFlex,
						},
					},
					Margin: linebot.FlexComponentMarginTypeXs,
					Flex:   &minFlex,
				},
			},
			BackgroundColor: "#3D4D4D",
			Margin:          linebot.FlexComponentMarginTypeNone,
		},
	}
}

// GenerateLineupPokemonsFlexComponent converts lineup pokemon to LINE flex message
func GenerateLineupPokemonsFlexComponent(allLineupPokemons []gd.LineupPokemon, seltectedSolt int) []linebot.FlexComponent {
	minFlex := 1
	withoutFlex := 0

	lineupPokemons := funk.Filter(allLineupPokemons, func(lineupPokemon gd.LineupPokemon) bool {
		return lineupPokemon.SlotNo == seltectedSolt
	}).([]gd.LineupPokemon)

	contents := func() []linebot.FlexComponent {
		rowContents := []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:    linebot.FlexComponentTypeText,
				Text:    fmt.Sprintf("順序 #%d", seltectedSolt),
				Size:    linebot.FlexTextSizeTypeMd,
				Align:   linebot.FlexComponentAlignTypeCenter,
				Gravity: linebot.FlexComponentGravityTypeCenter,
				Color:   "#FFFFFF",
				Flex:    &withoutFlex,
				Margin:  linebot.FlexComponentMarginTypeNone,
			},
			&linebot.SeparatorComponent{
				Type:   linebot.FlexComponentTypeSeparator,
				Color:  "#CDCDCD",
				Margin: linebot.FlexComponentMarginTypeMd,
			},
		}

		for _, lineupPokemon := range lineupPokemons {
			nameText := lineupPokemon.Name

			if lineupPokemon.Catchable {
				nameText = "☘️ " + lineupPokemon.Name
			}

			if lineupPokemon.ShinyAvailable {
				nameText += " ✨"
			}

			rowContents = append(
				rowContents,
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeVertical,
					Contents: []linebot.FlexComponent{
						&linebot.ImageComponent{
							Type:  linebot.FlexComponentTypeImage,
							Size:  "75px",
							URL:   lineupPokemon.ImageURL,
							Align: linebot.FlexComponentAlignTypeCenter,
						},
						&linebot.TextComponent{
							Type:    linebot.FlexComponentTypeText,
							Text:    nameText,
							Size:    linebot.FlexTextSizeTypeSm,
							Align:   linebot.FlexComponentAlignTypeCenter,
							Gravity: linebot.FlexComponentGravityTypeCenter,
							Color:   "#FFFFFF",
							Flex:    &minFlex,
							Margin:  linebot.FlexComponentMarginTypeNone,
						},
					},
					Margin: linebot.FlexComponentMarginTypeNone,
				},
			)
		}

		return rowContents
	}()

	return contents
}
