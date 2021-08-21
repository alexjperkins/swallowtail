package satoshi

import (
	"context"
	"fmt"
	"strings"
	"time"

	twitterclient "swallowtail/s.satoshi/clients"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/monzo/slog"
)

const (
	twitterConsumerID = "twitter-consumer"
)

var (
	// For ease of mocking purposes.
	defaultTwitterClient = twitterclient.New
)

type TwitterUserMetaData struct {
	Bio            string
	Name           string
	Id             string
	DiscordChannel string
	Emoji          string
	Twitter        string
	Twitch         string
	Youtube        string
	Tags           []string
	Filter         func(string) bool
}

func init() {
	registerSatoshiConsumer(twitterConsumerID, TwitterConsumer{
		Active: true,
	})
}

type TwitterConsumer struct {
	Active bool
}

func (tw TwitterConsumer) Receiver(ctx context.Context, c chan *SatoshiConsumerMessage, d chan struct{}, _ bool) {
	usersToConsumer := []string{}
	for _, user := range usernameMetadataMapping {
		usersToConsumer = append(usersToConsumer, user.Id)
	}

	cli := defaultTwitterClient()
	filterParams := &twitter.StreamFilterParams{
		Follow:        usersToConsumer,
		StallWarnings: twitter.Bool(true),
	}

	tweetHandler := postTweetToDiscordHandler(ctx, c, tw.Active)
	err := cli.NewStream(ctx, filterParams, tweetHandler)
	if err != nil {
		slog.Error(ctx, "Failed to start twitter stream.", map[string]string{
			"consumer_id": twitterConsumerID,
			"error":       err.Error(),
		})
		return
	}
	defer cli.StopStream()
	defer slog.Warn(ctx, "Stopping twitter client stream")
	slog.Info(ctx, "Stream created; waiting for tweets...")
	for {
		select {
		case <-ctx.Done():
		case <-d:
		}
	}
}

func (tw TwitterConsumer) IsActive() bool {
	return tw.Active
}

func postTweetToDiscordHandler(ctx context.Context, c chan<- *SatoshiConsumerMessage, isActive bool) func(tweet *twitter.Tweet) {
	return func(tweet *twitter.Tweet) {
		// We can skip RTs.
		if strings.HasPrefix(tweet.Text, "RT") {
			return
		}
		user, ok := getMetadataMapping(tweet.User.ScreenName)
		if !ok {
			return
		}
		content := formatTweetForDiscord(tweet)
		slog.Trace(context.TODO(), content)
		msg := &SatoshiConsumerMessage{
			ConsumerID:       twitterConsumerID,
			Message:          content,
			DiscordChannelID: user.DiscordChannel,
			Created:          time.Now(),
			IdempotencyKey:   fmt.Sprintf("%s-%s", tweet.User.ScreenName, tweet.CreatedAt),
			IsActive:         isActive,
			Metadata: map[string]string{
				"username":        user.Name,
				"tweet_timestamp": tweet.CreatedAt,
			},
		}
		select {
		case c <- msg:
		case <-ctx.Done():
		default:
			slog.Warn(ctx, "Dropping twitter msg; satoshi consumer channel blocked: %v", msg)
		}
	}
}

func formatTweetForDiscord(tweet *twitter.Tweet) string {
	return fmt.Sprintf("@%s [%v]: %s", tweet.User.ScreenName, tweet.CreatedAt, tweet.Text)
}
