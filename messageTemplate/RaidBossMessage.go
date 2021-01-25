package messageTemplate

import (
	"fmt"
	"strings"

	gd "pmgo-professor-willow/lineChatbot/gamedata"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

// GenerateRaidBossMessages converts raid bosses to LINE flex messages
func GenerateRaidBossMessages(raidBosses []gd.RaidBoss, raidTier string) []linebot.SendingMessage {
	bossChunks := funk.Chunk(raidBosses, 10).([][]gd.RaidBoss)

	return funk.Map(bossChunks, func(bossChunk []gd.RaidBoss) linebot.SendingMessage {
		return linebot.NewFlexMessage(
			fmt.Sprintf("%s 星團體戰列表", raidTier),
			&linebot.CarouselContainer{
				Type: linebot.FlexContainerTypeCarousel,
				Contents: funk.Map(
					bossChunk,
					GenerateRaidBossBubbleMessage,
				).([]*linebot.BubbleContainer),
			},
		)
	}).([]linebot.SendingMessage)
}

// GenerateRaidBossBubbleMessage converts raid boss to LINE bubble message
func GenerateRaidBossBubbleMessage(raidBoss gd.RaidBoss) *linebot.BubbleContainer {
	maxFlex := 10
	minFlex := 1
	withoutFlex := 0

	nameContents := func() []linebot.FlexComponent {
		title := raidBoss.Name

		if raidBoss.ShinyAvailable {
			title += " ✨"
		}

		results := []linebot.FlexComponent{
			&linebot.TextComponent{
				Type: linebot.FlexComponentTypeText,
				Text: title,
				Size: linebot.FlexTextSizeTypeLg,
				Flex: &maxFlex,
			},
			&linebot.ImageComponent{
				Type:  linebot.FlexComponentTypeImage,
				Size:  "20px",
				URL:   raidBoss.TypeURLs[0],
				Align: linebot.FlexComponentAlignTypeEnd,
				Flex:  &minFlex,
			},
		}

		if len(raidBoss.TypeURLs) > 1 {
			results = append(results, &linebot.ImageComponent{
				Type:  linebot.FlexComponentTypeImage,
				Size:  "20px",
				URL:   raidBoss.TypeURLs[1],
				Align: linebot.FlexComponentAlignTypeEnd,
				Flex:  &minFlex,
			})
		}

		return results
	}()

	boostedCPContents := func() []linebot.FlexComponent {
		results := []linebot.FlexComponent{
			// Empty text for padding.
			&linebot.TextComponent{
				Type:  linebot.FlexComponentTypeText,
				Text:  "_",
				Size:  linebot.FlexTextSizeTypeLg,
				Flex:  &maxFlex,
				Color: "#FFFFFF",
			},
			&linebot.ImageComponent{
				Type:  linebot.FlexComponentTypeImage,
				Size:  "20px",
				URL:   raidBoss.BoostedWeatherURLs[0],
				Align: linebot.FlexComponentAlignTypeEnd,
				Flex:  &withoutFlex,
			},
		}

		if len(raidBoss.BoostedWeatherURLs) > 1 {
			results = append(results, &linebot.ImageComponent{
				Type:  linebot.FlexComponentTypeImage,
				Size:  "20px",
				URL:   raidBoss.BoostedWeatherURLs[1],
				Align: linebot.FlexComponentAlignTypeEnd,
				Flex:  &withoutFlex,
			})
		}

		results = append(
			results,
			// Empty text for padding.
			&linebot.TextComponent{
				Type:  linebot.FlexComponentTypeText,
				Text:  "_",
				Size:  linebot.FlexTextSizeTypeLg,
				Flex:  &withoutFlex,
				Color: "#FFFFFF",
			},
			&linebot.TextComponent{
				Type:  linebot.FlexComponentTypeText,
				Text:  fmt.Sprintf("CP: %d - %d", raidBoss.BoostedCP.Min, raidBoss.BoostedCP.Max),
				Color: "#6C757D",
				Flex:  &withoutFlex,
				Align: linebot.FlexComponentAlignTypeEnd,
			},
		)

		return results
	}()

	return &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Size: linebot.FlexBubbleSizeTypeKilo,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: []linebot.FlexComponent{
						&linebot.ImageComponent{
							Type: linebot.FlexComponentTypeImage,
							Size: linebot.FlexImageSizeTypeMd,
							URL:  strings.ReplaceAll(raidBoss.ImageURL, "//images.weserv.nl/?w=200&il&url=", "https://"),
						},
					},
				},
				&linebot.BoxComponent{
					Type:     linebot.FlexComponentTypeBox,
					Layout:   linebot.FlexBoxLayoutTypeHorizontal,
					Contents: nameContents,
				},
				&linebot.SeparatorComponent{
					Type:   linebot.FlexComponentTypeSeparator,
					Color:  "#CDCDCD",
					Margin: linebot.FlexComponentMarginTypeXs,
				},
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: []linebot.FlexComponent{
						// Empty text for padding.
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  "_",
							Size:  linebot.FlexTextSizeTypeLg,
							Flex:  &maxFlex,
							Color: "#FFFFFF",
						},
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  fmt.Sprintf("CP: %d - %d", raidBoss.CP.Min, raidBoss.CP.Max),
							Color: "#6C757D",
							Flex:  &withoutFlex,
							Align: linebot.FlexComponentAlignTypeEnd,
						},
					},
				},
				&linebot.BoxComponent{
					Type:     linebot.FlexComponentTypeBox,
					Layout:   linebot.FlexBoxLayoutTypeHorizontal,
					Contents: boostedCPContents,
				},
			},
		},
	}
}
