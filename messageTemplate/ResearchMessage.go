package messageTemplate

import (
	"fmt"

	gd "pmgo-professor-willow/lineChatbot/gamedata"
	"pmgo-professor-willow/lineChatbot/messageTemplate/utils"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

type Category string
type ResearchCollection map[Category][]gd.Research

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
			return Category(research.Category) == category && len(research.RewardPokemons) > 0
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
				Contents: funk.Map(categories, func(category Category) *linebot.BubbleContainer {
					return GenerateResearchBubbleMessage(category, collections[category])
				}).([]*linebot.BubbleContainer),
			},
		),
	}
}

// GenerateResearchBubbleMessage converts researches to LINE bubble message
func GenerateResearchBubbleMessage(category Category, researches []gd.Research) *linebot.BubbleContainer {
	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeGiga,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeText,
			Layout: linebot.FlexBoxLayoutTypeHorizontal,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   string(category),
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
					}

					return row
				}).([][]linebot.FlexComponent),
			).([]linebot.FlexComponent),
			BackgroundColor: "#3D4D4D",
			Margin:          linebot.FlexComponentMarginTypeNone,
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

		results := []linebot.FlexComponent{
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
