package consumers

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
	}{}
	// TODO: We want non-blocking test cases here; the difficulty lies in the fact that
	// we can't mock the discord session.
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				ctx = context.Background()
				c   = make(chan *ConsumerMessage, 1)
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
	}{}
	// TODO: We want non-blocking test cases here; the difficulty lies in the fact that
	// we can't mock the discord session.
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				ctx = context.Background()
				c   = make(chan *ConsumerMessage, 1)
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

func TestContainsAstekz1To10kChallenge(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                       string
		content                    string
		modUsername                string
		doesContain1To10kChallenge bool
	}{
		{
			name:        "does-contain-challenge",
			modUsername: "Astekz",
			content: `
			---NEW---
			Astekz [￼astekz]: 1k - 10k spot challenge

			1105 spent on DOT at 17.1 with a stop of 16.4`,
			doesContain1To10kChallenge: true,
		},
		{
			name: "contains-but-not-astekz",
			content: `
			---NEW---
			rego [￼rego]: 1k - 10k spot challenge

			1105 spent on DOT at 17.1 with a stop of 16.4`,
			doesContain1To10kChallenge: false,
		},
		{
			name:        "from-astekz-but-not-challenge",
			modUsername: "Astekz",
			content: `
			Astekz [:nazar_amulet:astekz]: lets just close dot here guys, like next to nothing loss, i think it goes to 16.4. 
			`,
			doesContain1To10kChallenge: false,
		},
		{
			name:        "another-example",
			modUsername: "Astekz",
			content: `
			Astekz [:moneybag:1k-10k]: 1k -10k spot call risky but i like it

			1500 spent on shib at 0.008050
			`,
			doesContain1To10kChallenge: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			didContain1To10kChallenge := containsAstekz1To10kChallenge(tt.modUsername, tt.content)
			assert.Equal(t, tt.doesContain1To10kChallenge, didContain1To10kChallenge)
		})
	}
}
