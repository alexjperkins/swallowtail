package client

import (
	"context"
	"sync"
)

var (
	client Web3Client
	mu     sync.Mutex
)

type Web3Client interface {
	SubscribePendingTransactions(ctx context.Context) (<-chan *PendingTransactionEvent, error)
	StopReceiver()
}

func Init(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()
	if client != nil {
		return nil
	}
	c, err := NewInfuraClient(ctx)
	if err != nil {
		return err
	}
	client = c
	return nil
}

func UseMock() {
	mu.Lock()
	defer mu.Unlock()
	client = &MockClient{}
}
