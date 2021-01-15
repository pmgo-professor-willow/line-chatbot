package functions

import (
	"fmt"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/thoas/go-funk"
)

// GenerateRaidBossMessages converts raid bosses to LINE flex messages
func GenerateRaidBossMessages(raidBosses []RaidBoss, raidTier string) []linebot.SendingMessage {
	bossChunks := funk.Chunk(raidBosses, 10).([][]RaidBoss)

	return funk.Map(bossChunks, func(bossChunk []RaidBoss) linebot.SendingMessage {
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
func GenerateRaidBossBubbleMessage(raidBoss RaidBoss) *linebot.BubbleContainer {
	maxFlex := 10
	minFlex := 1
	withoutFlex := 0

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
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: append(
						[]linebot.FlexComponent{
							&linebot.TextComponent{
								Type: linebot.FlexComponentTypeText,
								Text: raidBoss.Name,
								Size: linebot.FlexTextSizeTypeLg,
								Flex: &maxFlex,
							},
						},
						funk.Map(raidBoss.TypeURLs, func(typeURL string) *linebot.ImageComponent {
							return &linebot.ImageComponent{
								Type:  linebot.FlexComponentTypeImage,
								Size:  "20px",
								URL:   typeURL,
								Align: linebot.FlexComponentAlignTypeEnd,
								Flex:  &minFlex,
							}
						}).([]linebot.FlexComponent)...,
					),
				},
				&linebot.SpacerComponent{
					Type: linebot.FlexComponentTypeSeparator,
				},
				&linebot.BoxComponent{
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: []linebot.FlexComponent{
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
					Type:   linebot.FlexComponentTypeBox,
					Layout: linebot.FlexBoxLayoutTypeHorizontal,
					Contents: append(
						funk.Map(raidBoss.BoostedWeatherURLs, func(boostedWeatherURL string) *linebot.ImageComponent {
							return &linebot.ImageComponent{
								Type:  linebot.FlexComponentTypeImage,
								Size:  "20px",
								URL:   boostedWeatherURL,
								Align: linebot.FlexComponentAlignTypeEnd,
								Flex:  &withoutFlex,
							}
						}).([]linebot.FlexComponent),
						&linebot.TextComponent{
							Type:  linebot.FlexComponentTypeText,
							Text:  fmt.Sprintf("CP: %d - %d", raidBoss.BoostedCP.Min, raidBoss.BoostedCP.Max),
							Color: "#6C757D",
							Flex:  &withoutFlex,
							Align: linebot.FlexComponentAlignTypeEnd,
						},
					),
				},
			},
		},
	}
}
