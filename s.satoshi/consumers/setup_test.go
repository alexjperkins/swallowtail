package consumers

import (
	"os"
	"testing"

	binanceclient "swallowtail/s.binance/client"
	twitterclient "swallowtail/s.satoshi/clients"
)

func SetupTest(m *testing.M) {
	// Mock the default twitter client.
	defaultTwitterClient = func() twitterclient.TwitterClient { return &twitterclient.MockTwitterClient{} }
	// Mock the default Binance client.
	binanceclient.UseMock()
	os.Exit(m.Run())
}
