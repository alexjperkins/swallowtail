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
	}{}
	// TODO: We want non-blocking test cases here; the difficulty lies in the fact that
	// we can't mock the discord session.
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
	}{}
	// TODO: We want non-blocking test cases here; the difficulty lies in the fact that
	// we can't mock the discord session.
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

func TestContainsTicker(t *testing.T) {
	// Set our map of binance asset pairs to something we can test against.
	originalBinanceAssetPairs := binanceAssetPairs
	binanceAssetPairs = map[string]bool{
		"btc": true,
		"eth": true,
	}
	t.Cleanup(func() {
		binanceAssetPairs = originalBinanceAssetPairs
	})

	tests := []struct {
		name            string
		content         string
		doesTickerExist bool
	}{
		{
			name:            "ticker_does_not_exist",
			content:         "this is some content with not ticker at all.",
			doesTickerExist: false,
		},
		{
			name:            "contains_stable_coin",
			content:         "this contains a stablecoin e.g usd usdt usdc",
			doesTickerExist: false,
		},
		{
			name:            "contains_bluntz styled_ticker",
			content:         "this contains btc/usd & eth/usd tickers",
			doesTickerExist: true,
		},
		{
			name:            "contains_ticker",
			content:         "contains regular styled ticker with newline \n btc and a few 5.9 \n sol",
			doesTickerExist: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res := containsTicker(tt.content)
			assert.Equal(t, tt.doesTickerExist, res)
		})
	}
}
