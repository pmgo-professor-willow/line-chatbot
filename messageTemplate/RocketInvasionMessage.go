package messageTemplate

import (
	"fmt"

	gd "pmgo-professor-willow/lineChatbot/gamedata"
	"pmgo-professor-willow/lineChatbot/messageTemplate/utils"
	pokemonUtils "pmgo-professor-willow/lineChatbot/utils"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/thoas/go-funk"
)

// GenerateRocketInvasionListMessages sends LINE quick reply messages
func GenerateRocketInvasionListMessages() []linebot.SendingMessage {
	return []linebot.SendingMessage{
		linebot.NewTextMessage(
			"想知道火箭隊手下還是幹部的陣容呢？",
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

	hasQuote := rocketInvasion.Quote != ""
	titleText := ""
	defaultTitleGravity := linebot.FlexComponentGravityTypeCenter

	if rocketInvasion.IsSpecial {
		titleText = rocketInvasion.Category
	} else {
		titleText = fmt.Sprintf("手下 (%s)", rocketInvasion.Category)
	}

	if hasQuote {
		defaultTitleGravity = linebot.FlexComponentGravityTypeBottom
	}

	titleContent := &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeVertical,
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:    linebot.FlexComponentTypeText,
				Text:    titleText,
				Size:    linebot.FlexTextSizeTypeLg,
				Align:   linebot.FlexComponentAlignTypeStart,
				Gravity: defaultTitleGravity,
				Margin:  linebot.FlexComponentMarginTypeXs,
				Color:   "#FFFFFF",
				Weight:  linebot.FlexTextWeightTypeBold,
				Flex:    &minFlex,
			},
		},
		Flex: &maxFlex,
	}

	// Append quote if exist
	if hasQuote {
		titleContent.Contents = append(
			titleContent.Contents,
			&linebot.TextComponent{
				Type:    linebot.FlexComponentTypeText,
				Text:    rocketInvasion.Quote,
				Size:    linebot.FlexTextSizeTypeSm,
				Align:   linebot.FlexComponentAlignTypeStart,
				Gravity: linebot.FlexComponentGravityTypeTop,
				Margin:  linebot.FlexComponentMarginTypeXs,
				Color:   "#CDCDCD",
				Weight:  linebot.FlexTextWeightTypeBold,
				Wrap:    true,
				Flex:    &minFlex,
			},
		)
	}

	headerContents := []linebot.FlexComponent{
		// Character image
		&linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.SpacerComponent{
					Size: linebot.FlexSpacerSizeTypeMd,
				},
				&linebot.ImageComponent{
					Type:    linebot.FlexComponentTypeImage,
					Size:    "100px",
					URL:     rocketInvasion.CharacterImageURL,
					Align:   linebot.FlexComponentAlignTypeCenter,
					Gravity: linebot.FlexComponentGravityTypeBottom,
					Margin:  linebot.FlexComponentMarginTypeNone,
					Flex:    &minFlex,
				},
			},
			Margin: linebot.FlexComponentMarginTypeNone,
			Flex:   &minFlex,
		},
		// Character name and quote
		titleContent,
	}

	// Add recommendation button
	if rocketInvasion.IsSpecial {
		// Keep total is 4 (1+2+1)
		maxFlex = 2
		headerContents = append(
			// Character image, name and quote
			headerContents,
			// Defeat recommendation
			&linebot.ButtonComponent{
				Type:  linebot.FlexComponentTypeButton,
				Style: linebot.FlexButtonStyleTypeLink,
				Color: "#FFFFFF",
				Action: &linebot.PostbackAction{
					Label:       "🔰弱點",
					Data:        fmt.Sprintf("rocketInvasion=%s", rocketInvasion.Category),
					DisplayText: fmt.Sprintf("請告訴我如何克制 %s 的陣容", titleText),
				},
				Gravity: linebot.FlexComponentGravityTypeCenter,
				Flex:    &maxFlex,
			},
		)
	}

	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeGiga,
		Header: &linebot.BoxComponent{
			Type:            linebot.FlexComponentTypeBox,
			Layout:          linebot.FlexBoxLayoutTypeHorizontal,
			Contents:        headerContents,
			BackgroundColor: "#455F60",
			PaddingTop:      linebot.FlexComponentPaddingTypeNone,
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
						// Lineup slot #1
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
						// Lineup slot #2
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
						// Lineup slot #3
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
			PaddingAll:      linebot.FlexComponentPaddingTypeNone,
		},
	}
}

// GenerateLineupPokemonsFlexComponent converts lineup pokemon to LINE flex message
func GenerateLineupPokemonsFlexComponent(allLineupPokemons []gd.LineupPokemon, seltectedSolt int) []linebot.FlexComponent {
	maxFlex := 2
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
				Margin:  linebot.FlexComponentMarginTypeMd,
			},
			&linebot.SeparatorComponent{
				Type:   linebot.FlexComponentTypeSeparator,
				Color:  "#CDCDCD",
				Margin: linebot.FlexComponentMarginTypeMd,
			},
		}

		for _, lineupPokemon := range lineupPokemons {
			avatarContents := []linebot.FlexComponent{}

			// Left of avatar row
			if lineupPokemon.Catchable {
				avatarContents = append(
					avatarContents,
					&linebot.TextComponent{
						Type:    linebot.FlexComponentTypeText,
						Text:    "☘️",
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
					URL:     lineupPokemon.ImageURL,
					Align:   linebot.FlexComponentAlignTypeCenter,
					Gravity: linebot.FlexComponentGravityTypeBottom,
					Flex:    &maxFlex,
				},
			)

			// Right of avatar row
			if lineupPokemon.ShinyAvailable {
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

			rowContents = append(
				rowContents,
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
							Type:    linebot.FlexComponentTypeText,
							Text:    lineupPokemon.Name,
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

		// Add padding
		rowContents = append(
			rowContents,
			&linebot.SpacerComponent{
				Size: linebot.FlexSpacerSizeTypeMd,
			},
		)

		return rowContents
	}()

	return contents
}

// GenerateRocketInvasionWeaknessMessage converts rocket invasions to LINE flex messages
func GenerateRocketInvasionWeaknessMessage(rocketInvasion gd.RocketInvasion) []linebot.SendingMessage {
	return []linebot.SendingMessage{
		linebot.NewFlexMessage(
			fmt.Sprintf("火箭隊成員與幹部陣容"),
			&linebot.CarouselContainer{
				Type: linebot.FlexContainerTypeCarousel,
				Contents: []*linebot.BubbleContainer{
					GenerateRocketInvasionWeaknessBubbleMessage(rocketInvasion),
				},
			},
		),
	}
}

// GenerateRocketInvasionWeaknessBubbleMessage converts rocket invasion weakness to LINE bubble message
func GenerateRocketInvasionWeaknessBubbleMessage(rocketInvasion gd.RocketInvasion) *linebot.BubbleContainer {
	maxFlex := 3
	minFlex := 1

	titleText := "的陣容分析"
	if rocketInvasion.IsSpecial {
		titleText = fmt.Sprintf("%s %s", rocketInvasion.Category, titleText)
	} else {
		titleText = fmt.Sprintf("手下 (%s) %s", rocketInvasion.Category, titleText)
	}

	titleContent := &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeVertical,
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:    linebot.FlexComponentTypeText,
				Text:    titleText,
				Size:    linebot.FlexTextSizeTypeLg,
				Align:   linebot.FlexComponentAlignTypeStart,
				Gravity: linebot.FlexComponentGravityTypeCenter,
				Margin:  linebot.FlexComponentMarginTypeXs,
				Color:   "#FFFFFF",
				Weight:  linebot.FlexTextWeightTypeBold,
				Flex:    &minFlex,
			},
		},
		Flex: &maxFlex,
	}

	headerContents := []linebot.FlexComponent{
		// Character image
		&linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.SpacerComponent{
					Size: linebot.FlexSpacerSizeTypeMd,
				},
				&linebot.ImageComponent{
					Type:    linebot.FlexComponentTypeImage,
					Size:    "100px",
					URL:     rocketInvasion.CharacterImageURL,
					Align:   linebot.FlexComponentAlignTypeCenter,
					Gravity: linebot.FlexComponentGravityTypeBottom,
					Margin:  linebot.FlexComponentMarginTypeNone,
					Flex:    &minFlex,
				},
			},
			Margin: linebot.FlexComponentMarginTypeNone,
			Flex:   &minFlex,
		},
		// Character name
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
			PaddingTop:      linebot.FlexComponentPaddingTypeNone,
			PaddingBottom:   linebot.FlexComponentPaddingTypeNone,
		},
		Body: &linebot.BoxComponent{
			Type:            linebot.FlexComponentTypeBox,
			Layout:          linebot.FlexBoxLayoutTypeVertical,
			Contents:        GenerateLineupPokemonsWeaknessFlexComponent(rocketInvasion.LineupPokemons),
			BackgroundColor: "#3D4D4D",
			Margin:          linebot.FlexComponentMarginTypeNone,
			PaddingAll:      linebot.FlexComponentPaddingTypeNone,
		},
	}
}

// GenerateLineupPokemonsWeaknessFlexComponent converts lineup pokemon weakness to LINE flex message
func GenerateLineupPokemonsWeaknessFlexComponent(lineupPokemons []gd.LineupPokemon) []linebot.FlexComponent {
	maxFlex := 3
	minFlex := 1
	withoutFlex := 0

	contents := func() []linebot.FlexComponent {
		rowContents := []linebot.FlexComponent{}

		for i, lineupPokemon := range lineupPokemons {
			contents := []linebot.FlexComponent{}

			superEffectiveTypes := pokemonUtils.GetWeaknessTypes(lineupPokemon.Types, 2.56, 10)
			effectiveTypes := pokemonUtils.GetWeaknessTypes(lineupPokemon.Types, 1.6, 2.56)
			notEffectiveTypes := pokemonUtils.GetWeaknessTypes(lineupPokemon.Types, 0.625, 0.999)
			superNotEffectiveTypes := pokemonUtils.GetWeaknessTypes(lineupPokemon.Types, 0, 0.624)

			if len(superEffectiveTypes) > 0 {
				contents = append(
					contents,
					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeHorizontal,
						Contents: append(
							[]linebot.FlexComponent{
								&linebot.TextComponent{
									Type:    linebot.FlexComponentTypeText,
									Text:    "極佳",
									Size:    linebot.FlexTextSizeTypeSm,
									Align:   linebot.FlexComponentAlignTypeStart,
									Gravity: linebot.FlexComponentGravityTypeCenter,
									Color:   "#CDF68B",
									Flex:    &withoutFlex,
									Margin:  linebot.FlexComponentMarginTypeNone,
								},
							},
							funk.Map(
								superEffectiveTypes,
								utils.GetIcon,
							).([]linebot.FlexComponent)...,
						),
						Flex:   &minFlex,
						Margin: linebot.FlexComponentMarginTypeSm,
					},
				)
			}

			if len(effectiveTypes) > 0 {
				contents = append(
					contents,
					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeHorizontal,
						Contents: append(
							[]linebot.FlexComponent{
								&linebot.TextComponent{
									Type:    linebot.FlexComponentTypeText,
									Text:    "絕佳",
									Size:    linebot.FlexTextSizeTypeSm,
									Align:   linebot.FlexComponentAlignTypeStart,
									Gravity: linebot.FlexComponentGravityTypeCenter,
									Color:   "#CDF68B",
									Flex:    &withoutFlex,
									Margin:  linebot.FlexComponentMarginTypeNone,
								},
							},
							funk.Map(
								effectiveTypes,
								utils.GetIcon,
							).([]linebot.FlexComponent)...,
						),
						Flex:   &minFlex,
						Margin: linebot.FlexComponentMarginTypeSm,
					},
				)
			}

			contents = append(
				contents,
				&linebot.SeparatorComponent{
					Type:   linebot.FlexComponentTypeSeparator,
					Color:  "#CDCDCD",
					Margin: linebot.FlexComponentMarginTypeSm,
				},
			)

			if len(notEffectiveTypes) > 0 {
				contents = append(
					contents,
					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeHorizontal,
						Contents: append(
							[]linebot.FlexComponent{
								&linebot.TextComponent{
									Type:    linebot.FlexComponentTypeText,
									Text:    "不好",
									Size:    linebot.FlexTextSizeTypeSm,
									Align:   linebot.FlexComponentAlignTypeStart,
									Gravity: linebot.FlexComponentGravityTypeCenter,
									Color:   "#CDCDCD",
									Flex:    &withoutFlex,
									Margin:  linebot.FlexComponentMarginTypeNone,
								},
							},
							funk.Map(
								notEffectiveTypes,
								utils.GetIcon,
							).([]linebot.FlexComponent)...,
						),
						Flex:   &minFlex,
						Margin: linebot.FlexComponentMarginTypeSm,
					},
				)
			}

			if len(superNotEffectiveTypes) > 0 {
				contents = append(
					contents,
					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeHorizontal,
						Contents: append(
							[]linebot.FlexComponent{
								&linebot.TextComponent{
									Type:    linebot.FlexComponentTypeText,
									Text:    "極差",
									Size:    linebot.FlexTextSizeTypeSm,
									Align:   linebot.FlexComponentAlignTypeStart,
									Gravity: linebot.FlexComponentGravityTypeCenter,
									Color:   "#CDCDCD",
									Flex:    &withoutFlex,
									Margin:  linebot.FlexComponentMarginTypeNone,
								},
							},
							funk.Map(
								superNotEffectiveTypes,
								utils.GetIcon,
							).([]linebot.FlexComponent)...,
						),
						Flex:   &minFlex,
						Margin: linebot.FlexComponentMarginTypeSm,
					},
				)
			}

			// Add padding
			contents = append(
				contents,
				&linebot.SpacerComponent{
					Size: linebot.FlexSpacerSizeTypeSm,
				},
			)

			rowContents = append(
				rowContents,
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: []linebot.FlexComponent{
						&linebot.SpacerComponent{
							Size: linebot.FlexSpacerSizeTypeSm,
						},
						// Slot No
						&linebot.TextComponent{
							Type:    linebot.FlexComponentTypeText,
							Text:    fmt.Sprintf("#%d", lineupPokemon.SlotNo),
							Size:    linebot.FlexTextSizeTypeSm,
							Align:   linebot.FlexComponentAlignTypeStart,
							Gravity: linebot.FlexComponentGravityTypeCenter,
							Color:   "#FFFFFF",
							Flex:    &withoutFlex,
						},
						// Image
						&linebot.ImageComponent{
							Type:    linebot.FlexComponentTypeImage,
							Size:    "75px",
							URL:     lineupPokemon.ImageURL,
							Align:   linebot.FlexComponentAlignTypeStart,
							Gravity: linebot.FlexComponentGravityTypeCenter,
							Flex:    &withoutFlex,
						},
						// Name and types
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Contents: []linebot.FlexComponent{
								// Name
								&linebot.TextComponent{
									Type:    linebot.FlexComponentTypeText,
									Text:    lineupPokemon.Name,
									Size:    linebot.FlexTextSizeTypeSm,
									Align:   linebot.FlexComponentAlignTypeStart,
									Gravity: linebot.FlexComponentGravityTypeBottom,
									Color:   "#FFFFFF",
									Flex:    &minFlex,
									Margin:  linebot.FlexComponentMarginTypeNone,
									Wrap:    true,
								},
								// Types
								&linebot.BoxComponent{
									Type:     linebot.FlexComponentTypeBox,
									Layout:   linebot.FlexBoxLayoutTypeHorizontal,
									Contents: funk.Map(lineupPokemon.Types, utils.GetIcon).([]linebot.FlexComponent),
									Flex:     &minFlex,
									Margin:   linebot.FlexComponentMarginTypeNone,
								},
								&linebot.SpacerComponent{
									Size: linebot.FlexSpacerSizeTypeMd,
								},
							},
						},
						// Reason
						&linebot.BoxComponent{
							Type:     linebot.FlexComponentTypeBox,
							Layout:   linebot.FlexBoxLayoutTypeVertical,
							Contents: contents,
							Flex:     &maxFlex,
							Margin:   linebot.FlexComponentMarginTypeXs,
						},
					},
					Margin: linebot.FlexComponentMarginTypeNone,
				},
			)

			if i+1 != len(lineupPokemons) {
				rowContents = append(
					rowContents,
					&linebot.SeparatorComponent{
						Type:   linebot.FlexComponentTypeSeparator,
						Color:  "#CDCDCD",
						Margin: linebot.FlexComponentMarginTypeNone,
					},
				)
			}
		}

		// Add padding
		rowContents = append(
			rowContents,
			&linebot.SpacerComponent{
				Size: linebot.FlexSpacerSizeTypeMd,
			},
		)

		return rowContents
	}()

	return contents
}
