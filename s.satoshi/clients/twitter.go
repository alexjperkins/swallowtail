package clients

import (
	"context"
	"fmt"
	"swallowtail/libraries/util"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/monzo/slog"
	"github.com/monzo/terrors"
)

var (
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
	bearerToken       string
)

var (
	defaultTwitterClient TwitterClient
)

func init() {
	consumerKey = util.SetEnv("TWITTER_API_KEY")
	consumerSecret = util.SetEnv("TWITTER_API_SECRET")

	accessToken = util.SetEnv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = util.SetEnv("TWITTER_ACCESS_TOKEN_SECRET")

	bearerToken = util.SetEnv("TWITTER_BEARER_TOKEN")
}

// TwitterClient interface
type TwitterClient interface {
	NewStream(ctx context.Context, filter *twitter.StreamFilterParams, handler func(tweet *twitter.Tweet)) error
	StopStream() bool
}

type twitterClient struct {
	c      *twitter.Client
	stream *twitter.Stream
}

func (tc *twitterClient) NewStream(ctx context.Context, filter *twitter.StreamFilterParams, handler func(tweet *twitter.Tweet)) error {
	s, err := tc.c.Streams.Filter(filter)
	if err != nil {
		return terrors.Augment(err, "Failed to create new twitter stream", nil)
	}
	tc.stream = s

	d := twitter.NewSwitchDemux()
	d.Tweet = handler
	go d.Handle(tc.stream.Messages)
	return nil
}

func (tc *twitterClient) StopStream() bool {
	if tc.stream == nil {
		return false
	}
	tc.stream.Stop()
	return true
}

func New() TwitterClient {
	if defaultTwitterClient != nil {
		return defaultTwitterClient
	}
	cfg := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := cfg.Client(oauth1.NoContext, token)
	cli := twitter.NewClient(httpClient)

	// Verify that the credentials are indeed correct
	user, rsp, err := cli.Accounts.VerifyCredentials(
		&twitter.AccountVerifyParams{},
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to twitter: %v %v", rsp, err))
	}

	slog.Info(nil, "Twitter bot connected: %v %v", user, rsp)
	defaultTwitterClient = &twitterClient{
		c: cli,
	}
	return defaultTwitterClient
}
