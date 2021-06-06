package satoshi

import (
	"testing"

	discordproto "swallowtail/s.discord/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublishATHMSG(t *testing.T) {
	t.Parallel()

	var (
		ch           = make(chan *SatoshiConsumerMessage, 1)
		msg          = "test-msg"
		symbol       = "BTC"
		currentPrice = 10.0
		currentATH   = 12.0
	)

	// Publish first attempt & consumer directly after.
	publishATHMsg(ch, msg, symbol, currentPrice, currentATH)

	e := <-ch
	require.NotNil(t, e)
	assert.Equal(t, msg, e.Message)
	assert.Equal(t, discordproto.DiscordSatoshiAlertsChannel, e.DiscordChannelID)

	// Publish two attemmpts; the second attempt to publish to a channel should be blocked,
	// meaning we log and ignore.
	publishATHMsg(ch, msg, symbol, currentPrice, currentATH)
	publishATHMsg(ch, "This is the wrong msg", "DOGE", currentPrice, currentATH)

	// Consumer and assert the message is correct.
	e = <-ch
	require.NotNil(t, e)
	assert.Equal(t, msg, e.Message)
	assert.Equal(t, discordproto.DiscordSatoshiAlertsChannel, e.DiscordChannelID)

	// Try once more to confirm that we don't see the "lost" message.
	publishATHMsg(ch, msg, symbol, currentPrice, currentATH)

	e = <-ch
	require.NotNil(t, e)
	assert.Equal(t, msg, e.Message)
	assert.Equal(t, discordproto.DiscordSatoshiAlertsChannel, e.DiscordChannelID)
}
