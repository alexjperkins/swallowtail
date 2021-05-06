package clients

import (
	"fmt"
	"swallowtail/libraries/util"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/monzo/slog"
)

var (
	consumerKey       string
	consumerSecret    string
	accessToken       string
	accessTokenSecret string
	bearerToken       string

	mtx sync.Mutex
)

type TwitterClient struct {
	Client *twitter.Client
}

func init() {
	mtx.Lock()
	defer mtx.Unlock()

	consumerKey = util.SetEnv("TWITTER_API_KEY")
	consumerSecret = util.SetEnv("TWITTER_API_SECRET")

	accessToken = util.SetEnv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = util.SetEnv("TWITTER_ACCESS_TOKEN_SECRET")

	bearerToken = util.SetEnv("TWITTER_BEARER_TOKEN")
}

func New() *TwitterClient {
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

	return &TwitterClient{
		Client: cli,
	}
}
