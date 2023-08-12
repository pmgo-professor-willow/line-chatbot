package messageTemplate

import (
	"fmt"

	gd "pmgo-professor-willow/lineChatbot/gamedata"
	"pmgo-professor-willow/lineChatbot/messageTemplate/utils"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/thoas/go-funk"
)

type Category string
type ResearchCollection map[Category][]gd.Research

// GenerateResearchListMessages sends LINE quick reply messages
func GenerateResearchListMessages() []linebot.SendingMessage {
	return []linebot.SendingMessage{
		linebot.NewTextMessage(
			"想知道哪一類的田野調查課題呢？",
		).WithQuickReplies(
			linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"https://raw.githubusercontent.com/pmgo-professor-willow/line-chatbot/main/assets/researches/event.png",
					&linebot.PostbackAction{
						Label:       "活動限定",
						Data:        "research=event",
						DisplayText: "請列出「活動限定」田野調查。",
					},
				),
				linebot.NewQuickReplyButton(
					"https://raw.githubusercontent.com/pmgo-professor-willow/line-chatbot/main/assets/researches/catching_and_throwing.png",
					&linebot.PostbackAction{
						Label:       "捕捉與投球",
						Data:        "research=catching_and_throwing",
						DisplayText: "請列出「捕捉與投球」相關田野調查。",
					},
				),
				linebot.NewQuickReplyButton(
					"https://raw.githubusercontent.com/pmgo-professor-willow/line-chatbot/main/assets/researches/rocket.png",
					&linebot.PostbackAction{
						Label:       "GO 火箭隊",
						Data:        "research=rocket",
						DisplayText: "請列出「GO 火箭隊」相關田野調查。",
					},
				),
				linebot.NewQuickReplyButton(
					"",
					&linebot.PostbackAction{
						Label:       "其它田野調查",
						Data:        "research=others",
						DisplayText: "請列出其它田野調查。",
					},
				),
			),
		),
	}
}

// GenerateResearchMessages converts researches to LINE flex messages
func GenerateResearchMessages(researches []gd.Research) []linebot.SendingMessage {
	if utils.IsEmpty(researches) {
		return utils.GenerateEmptyReasonMessage()
	}

	categories := funk.Uniq(funk.Map(researches, func(research gd.Research) Category {
		return Category(research.Category)
	})).([]Category)

	collections := make(ResearchCollection)
	// Append into collections.
	funk.ForEach(categories, func(category Category) {
		collections[category] = funk.Filter(researches, func(research gd.Research) bool {
			// Filter RewardPokemons or RewardPokemonMegaCandies is not empty.
			return Category(research.Category) == category && (len(research.RewardPokemons) > 0 || len(research.RewardPokemonMegaCandies) > 0)
		}).([]gd.Research)
	})
	// Ignore empty collection.
	categories = funk.Filter(categories, func(category Category) bool {
		return len(collections[category]) > 0
	}).([]Category)

	return []linebot.SendingMessage{
		linebot.NewFlexMessage(
			"田野調查課題獎勵一覽",
			&linebot.CarouselContainer{
				Type: linebot.FlexContainerTypeCarousel,
				Contents: funk.FlattenDeep(
					funk.Map(categories, func(category Category) []*linebot.BubbleContainer {
						recommandMaxRow := 9
						totalRowLength := 0
						researchChunks := [][]gd.Research{}

						funk.ForEach(collections[category], func(research gd.Research) {
							rowLength := len(research.RewardPokemons) + len(research.RewardPokemonMegaCandies)
							isFirstResearch := len(researchChunks) == 0 && totalRowLength == 0
							isOutOfRecommandation := totalRowLength+rowLength > recommandMaxRow

							if isFirstResearch || isOutOfRecommandation {
								researchChunks = append(
									researchChunks,
									[]gd.Research{
										research,
									},
								)
								totalRowLength = rowLength
							} else {
								researchChunks[len(researchChunks)-1] = append(
									researchChunks[len(researchChunks)-1],
									research,
								)
								totalRowLength += rowLength
							}
						})

						return funk.Map(researchChunks, func(researches []gd.Research) *linebot.BubbleContainer {
							index := funk.IndexOf(researchChunks, researches)
							title := string(category)
							if index > 0 {
								title = fmt.Sprintf("%s-%d", title, index+1)
							}

							return GenerateResearchBubbleMessage(title, researches)
						}).([]*linebot.BubbleContainer)
					}).([][]*linebot.BubbleContainer),
				).([]*linebot.BubbleContainer),
			},
		),
	}
}

// GenerateResearchBubbleMessage converts researches to LINE bubble message
func GenerateResearchBubbleMessage(title string, researches []gd.Research) *linebot.BubbleContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeGiga,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeText,
			Layout: linebot.FlexBoxLayoutTypeHorizontal,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   title,
					Size:   linebot.FlexTextSizeTypeLg,
					Align:  linebot.FlexComponentAlignTypeEnd,
					Margin: linebot.FlexComponentMarginTypeNone,
					Color:  "#FFFFFF",
					Weight: linebot.FlexTextWeightTypeBold,
				},
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   " 調查課題",
					Size:   linebot.FlexTextSizeTypeLg,
					Align:  linebot.FlexComponentAlignTypeStart,
					Margin: linebot.FlexComponentMarginTypeNone,
					Color:  "#FFFFFF",
				},
			},
			BackgroundColor: "#455F60",
		},
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: funk.FlattenDeep(
				funk.Map(researches, func(research gd.Research) []linebot.FlexComponent {
					row := []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:     linebot.FlexComponentTypeBox,
							Layout:   linebot.FlexBoxLayoutTypeHorizontal,
							Contents: GenerateResearchFlexComponent(research),
							Margin:   linebot.FlexComponentMarginTypeXs,
						},
					}

					// Append separator if not last row.
					if funk.IndexOf(researches, research) != len(researches)-1 {
						row = append(
							row,
							&linebot.SeparatorComponent{
								Color: "#DCDCDC",
							},
						)
					} else {
						// Add a spacer for last row
						row = append(
							row,
							&linebot.SpacerComponent{
								Size: linebot.FlexSpacerSizeTypeMd,
							},
						)
					}

					return row
				}).([][]linebot.FlexComponent),
			).([]linebot.FlexComponent),
			BackgroundColor: "#3D4D4D",
			Margin:          linebot.FlexComponentMarginTypeNone,
			PaddingAll:      linebot.FlexComponentPaddingTypeNone,
			PaddingStart:    linebot.FlexComponentPaddingTypeMd,
			PaddingEnd:      linebot.FlexComponentPaddingTypeMd,
		},
	}
}

// GenerateResearchFlexComponent converts researches to LINE bubble message
func GenerateResearchFlexComponent(research gd.Research) []linebot.FlexComponent {
	maxFlex := 5
	minFlex := 4
	withoutFlex := 0

	contents := func() []linebot.FlexComponent {
		pokemonAvatarContents := []linebot.FlexComponent{}

		// Rewawrd pokemons.
		for _, rewardPokemon := range research.RewardPokemons {
			shinyAvailableText := " "

			if rewardPokemon.ShinyAvailable {
				shinyAvailableText += "✨"
			}

			pokemonAvatarContents = append(
				pokemonAvatarContents,
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: []linebot.FlexComponent{
						&linebot.ImageComponent{
							Type:  linebot.FlexComponentTypeImage,
							Size:  "75px",
							URL:   rewardPokemon.ImageURL,
							Align: linebot.FlexComponentAlignTypeCenter,
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.TextComponent{
									Type:    linebot.FlexComponentTypeText,
									Text:    shinyAvailableText,
									Size:    linebot.FlexTextSizeTypeSm,
									Align:   linebot.FlexComponentAlignTypeStart,
									Gravity: linebot.FlexComponentGravityTypeCenter,
									Color:   "#FFFFFF",
									Flex:    &withoutFlex,
								},
								&linebot.TextComponent{
									Type:    linebot.FlexComponentTypeText,
									Text:    fmt.Sprintf("CP %d", rewardPokemon.CP.Max),
									Size:    linebot.FlexTextSizeTypeSm,
									Align:   linebot.FlexComponentAlignTypeEnd,
									Gravity: linebot.FlexComponentGravityTypeCenter,
									Color:   "#FFFFFF",
									Flex:    &maxFlex,
								},
							},
							Margin: linebot.FlexComponentMarginTypeNone,
						},
					},
					Margin: linebot.FlexComponentMarginTypeNone,
				},
			)

			// Append separator if not last row.
			if funk.IndexOf(research.RewardPokemons, rewardPokemon) != len(research.RewardPokemons)-1 {
				pokemonAvatarContents = append(
					pokemonAvatarContents,
					&linebot.SeparatorComponent{
						Color: "#ECECEC",
					},
				)
			}
		}

		// Rewawrd mega candies of pokemons.
		for _, rewardPokemonMegaCandy := range research.RewardPokemonMegaCandies {
			pokemonAvatarContents = append(
				pokemonAvatarContents,
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: []linebot.FlexComponent{
						&linebot.ImageComponent{
							Type:  linebot.FlexComponentTypeImage,
							Size:  "75px",
							URL:   rewardPokemonMegaCandy.ImageURL,
							Align: linebot.FlexComponentAlignTypeCenter,
						},
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeHorizontal,
							Contents: []linebot.FlexComponent{
								&linebot.ImageComponent{
									Type:    linebot.FlexComponentTypeImage,
									Size:    "25px",
									URL:     rewardPokemonMegaCandy.MegaCandyImageUrl,
									Align:   linebot.FlexComponentAlignTypeCenter,
									Gravity: linebot.FlexComponentGravityTypeCenter,
								},
								&linebot.TextComponent{
									Type:    linebot.FlexComponentTypeText,
									Text:    fmt.Sprintf("x %d", rewardPokemonMegaCandy.Count),
									Size:    linebot.FlexTextSizeTypeMd,
									Align:   linebot.FlexComponentAlignTypeStart,
									Gravity: linebot.FlexComponentGravityTypeCenter,
									Color:   "#FFFFFF",
									Margin:  linebot.FlexComponentMarginTypeSm,
								},
							},
						},
					},
					Margin: linebot.FlexComponentMarginTypeNone,
				},
			)

			// Append separator if not last row.
			if funk.IndexOf(research.RewardPokemonMegaCandies, rewardPokemonMegaCandy) != len(research.RewardPokemonMegaCandies)-1 {
				pokemonAvatarContents = append(
					pokemonAvatarContents,
					&linebot.SeparatorComponent{
						Color: "#ECECEC",
					},
				)
			}
		}

		results := []linebot.FlexComponent{
			// research.Description
			&linebot.TextComponent{
				Type:    linebot.FlexComponentTypeText,
				Text:    research.Description,
				Size:    linebot.FlexTextSizeTypeMd,
				Align:   linebot.FlexComponentAlignTypeStart,
				Flex:    &minFlex,
				Wrap:    true,
				Gravity: linebot.FlexComponentGravityTypeCenter,
				Color:   "#FFFFFF",
				Margin:  linebot.FlexComponentMarginTypeSm,
			},
			// pokemonAvatarContents
			&linebot.BoxComponent{
				Type:     linebot.FlexComponentTypeBox,
				Layout:   linebot.FlexBoxLayoutTypeVertical,
				Contents: pokemonAvatarContents,
				Flex:     &maxFlex,
				Margin:   linebot.FlexComponentMarginTypeNone,
			},
		}

		return results
	}()

	return contents
}
