package client

import (
	"context"
)

type MockClient struct{}

func (m *MockClient) SubscribePendingTransactions(ctx context.Context) (<-chan *PendingTransactionEvent, error) {
	return nil, nil
}

func (m *MockClient) StopReceiver() {}
