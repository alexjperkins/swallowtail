package satoshi

import (
	"context"
	discordproto "swallowtail/s.discord/proto"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleModMessages(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		isActive      bool
		messageCreate *discordgo.MessageCreate
		shouldPublish bool
	}{
		{
			name:     "inactive-and-should-publish",
			isActive: false,
			messageCreate: &discordgo.MessageCreate{
				Message: &discordgo.Message{
					Attachments: []*discordgo.MessageAttachment{
						{
							URL: "this-is-a-url",
						},
					},
					ChannelID: discordproto.DiscordMoonModMessagesChannel,
					Content:   "Harry Potter and the Chamber of Secrets",
				},
			},
			shouldPublish: true,
		},
		{
			name:     "active-and-should-publish",
			isActive: true,
			messageCreate: &discordgo.MessageCreate{
				Message: &discordgo.Message{
					Attachments: []*discordgo.MessageAttachment{
						{
							URL: "this-is-a-url-2",
						},
					},
					ChannelID: discordproto.DiscordMoonModMessagesChannel,
					Content:   "Harry Potter and the Goblet of Fire",
				},
			},
			shouldPublish: true,
		},
		{
			name:     "active-but-should-not-publish",
			isActive: false,
			messageCreate: &discordgo.MessageCreate{
				Message: &discordgo.Message{
					Attachments: []*discordgo.MessageAttachment{
						{
							URL: "this-is-a-url",
						},
					},
					ChannelID: "Invalid channel id",
					Content:   "Harry Potter and the Chamber of Secrets",
				},
			},
			shouldPublish: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				ctx = context.Background()
				c   = make(chan *SatoshiConsumerMessage, 1)
			)

			// Publish message twice; we shouldn't get blocked on the second time.
			f := handleModMessages(ctx, c, tt.isActive)
			f(nil, tt.messageCreate)
			f(nil, tt.messageCreate)

			// Run assertions
			switch {
			case !tt.shouldPublish:
				assert.Len(t, c, 0)
			default:
				// Attempt to consume.
				e := <-c
				require.NotNil(t, e)
				assert.Equal(t, tt.isActive, e.IsActive)
				assert.Equal(t, tt.messageCreate.Content, e.Message)
				assert.Equal(t, tt.messageCreate.Attachments, e.Attachments)
				assert.Equal(t, discordproto.DiscordSatoshiModMessagesChannel, e.DiscordChannelID)
			}
		})
	}
}

func TestHandleSwingMessages(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		isActive      bool
		messageCreate *discordgo.MessageCreate
		shouldPublish bool
	}{
		{
			name:     "inactive-and-should-publish",
			isActive: false,
			messageCreate: &discordgo.MessageCreate{
				Message: &discordgo.Message{
					Attachments: []*discordgo.MessageAttachment{
						{
							URL: "this-is-a-url",
						},
					},
					ChannelID: discordproto.DiscordMoonSwingGroupChannel,
					Content:   "Harry Potter and the Chamber of Secrets",
				},
			},
			shouldPublish: true,
		},
		{
			name:     "active-and-should-publish",
			isActive: true,
			messageCreate: &discordgo.MessageCreate{
				Message: &discordgo.Message{
					Attachments: []*discordgo.MessageAttachment{
						{
							URL: "this-is-a-url-2",
						},
					},
					ChannelID: discordproto.DiscordMoonSwingGroupChannel,
					Content:   "Harry Potter and the Goblet of Fire",
				},
			},
			shouldPublish: true,
		},
		{
			name:     "active-but-should-not-publish",
			isActive: false,
			messageCreate: &discordgo.MessageCreate{
				Message: &discordgo.Message{
					Attachments: []*discordgo.MessageAttachment{
						{
							URL: "this-is-a-url",
						},
					},
					ChannelID: "Invalid channel id",
					Content:   "Harry Potter and the Chamber of Secrets",
				},
			},
			shouldPublish: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				ctx = context.Background()
				c   = make(chan *SatoshiConsumerMessage, 1)
			)

			// Publish message twice; we shouldn't get blocked on the second time.
			f := handleSwingMessages(ctx, c, tt.isActive)
			f(nil, tt.messageCreate)
			f(nil, tt.messageCreate)

			// Run assertions
			switch {
			case !tt.shouldPublish:
				assert.Len(t, c, 0)
			default:
				// Attempt to consume.
				e := <-c
				require.NotNil(t, e)

				assert.Equal(t, tt.isActive, e.IsActive)
				assert.Equal(t, tt.messageCreate.Content, e.Message)
				assert.Equal(t, tt.messageCreate.Attachments, e.Attachments)
				assert.Equal(t, discordproto.DiscordSatoshiSwingsChannel, e.DiscordChannelID)
			}
		})
	}
}
