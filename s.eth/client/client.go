package client

import (
	"context"
	"swallowtail/s.eth/domain"
)

type ETHClient interface {
}

var (
	client EthClient
)

type EthClient interface {
	SubscribePendingTransactions(ctx context.Context) (<-chan *domain.EthMempoolTxEvent, error)
	StopReceiver()
}

func Init(ctx context.Context) error {
	c, err := NewInfuraClient(ctx)
	if err != nil {
		return err
	}
	client = c
	return nil
}

func UseMock() {
}
