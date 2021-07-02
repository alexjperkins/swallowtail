package clients

import (
	"context"

	"github.com/dghubble/go-twitter/twitter"
)

type MockTwitterClient struct {
}

func (m *MockTwitterClient) NewStream(ctx context.Context, filter *twitter.StreamFilterParams, handler func(tweet *twitter.Tweet)) error {
	return nil
}

func (m *MockTwitterClient) StopStream() bool {
	return true
}
