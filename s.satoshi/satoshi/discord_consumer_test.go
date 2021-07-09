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

func TestContains1To10kChallenge(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                       string
		content                    string
		doesContain1To10kChallenge bool
	}{
		{
			name: "does-contain-challenge",
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
			name: "from-astekz-but-not-challenge",
			content: `
			Astekz [:nazar_amulet:astekz]: lets just close dot here guys, like next to nothing loss, i think it goes to 16.4. 
			`,
			doesContain1To10kChallenge: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			didContain1To10kChallenge := contains1To10kChallenge(tt.content)
			assert.Equal(t, tt.doesContain1To10kChallenge, didContain1To10kChallenge)
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
		{
			name: "contains_ticker_and_stablecoin",
			content: `
			rego [￼rego]: BTCUSD - 1D - SWING LONG IDEA - BIG SIZE STRATEGY

			Since Bitcoin is at the bottom of our range, and funding has provided a negative rate
			for more than a few days now, it is time for me to build a swing long.

			We've already added our first position size of 15% at 33,700.
			Depending on the situation, we can expect to see lower. In this case, I want to add. 

			The point of this position is to hold it over a long period of time.
			Going underwater/sideways does not matter if funding is negative.
			With a clear invalidation, holding this position for a long period of time will pay
			us as long as our rate remains negative.

			I am looking to add at points of interest: 32,800, 31,600, and 30,000.
			I am currently not interested in longing higher. If the conditions change we will add accordingly.

			Here is my chart for reference:
			`,
			doesTickerExist: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := containsTicker(tt.content)
			assert.Equal(t, tt.doesTickerExist, res)
		})
	}
}
