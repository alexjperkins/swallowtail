package satoshi

import (
	"context"
	"fmt"
	"testing"

	discordproto "swallowtail/s.discord/proto"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatTweetForDiscord(t *testing.T) {
	t.Parallel()

	tests := []struct {
		tweet          *twitter.Tweet
		expectedOutput string
	}{
		{
			tweet: &twitter.Tweet{
				Text:      "this-is-a-tweet",
				CreatedAt: "created-at",
				User: &twitter.User{
					ScreenName: "HarryPotter",
				},
			},
			expectedOutput: "@HarryPotter [created-at]: this-is-a-tweet",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.tweet.Text, func(t *testing.T) {
			t.Parallel()

			res := formatTweetForDiscord(tt.tweet)
			assert.Equal(t, tt.expectedOutput, res)
		})
	}
}

func TestPostTweetToDiscordHandler(t *testing.T) {
	var (
		username = "test-username"
		text     = "A body of text"
	)
	tests := []struct {
		name            string
		tweet           *twitter.Tweet
		isActive        bool
		expectedMessage string
	}{
		{
			name: "active_and_should_publish",
			tweet: &twitter.Tweet{
				Text:      text,
				CreatedAt: "created-at",
				User: &twitter.User{
					ScreenName: username,
				},
			},
			expectedMessage: "@test-username [created-at]: A body of text",
		},
	}

	orig := usernameMetadataMapping
	usernameMetadataMapping = map[string]*TwitterUserMetaData{
		username: {
			DiscordChannel: discordproto.DiscordSatoshiTestingChannel,
		},
	}
	t.Cleanup(func() {
		usernameMetadataMapping = orig
	})

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				ctx = context.Background()
				ch  = make(chan *SatoshiConsumerMessage, 1)
			)
			handler := postTweetToDiscordHandler(ctx, ch, tt.isActive)
			// Post twice; we shouldn't block.
			handler(tt.tweet)
			handler(tt.tweet)

			require.Len(t, ch, 1)
			e := <-ch
			require.NotNil(t, e)

			assert.Equal(t, e.DiscordChannelID, discordproto.DiscordSatoshiTestingChannel)
			assert.Equal(t, e.IsActive, tt.isActive)
			assert.Equal(t, fmt.Sprintf("%s-%s", tt.tweet.User.ScreenName, tt.tweet.CreatedAt), e.IdempotencyKey)
			assert.Equal(t, tt.expectedMessage, e.Message)
		})
	}
}
