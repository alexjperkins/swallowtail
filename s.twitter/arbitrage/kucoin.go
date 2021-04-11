package arbitrage

import (
	"swallowtail/s.twitter/clients"
)

var (
	ArbitrageKucoinClientID = "arbitrage-kucion-client"
)

func init() {
	register(ArbitrageKucoinClientID,
		&ArbitrageKucoinClient{
			// TODO: thread default timeout
			i: &clients.KucoinClient{},
		},
	)
}

type ArbitrageKucoinClient struct {
	i *clients.KucoinClient
}

func (abc *ArbitrageKucoinClient) GetPrice(symbol string) float64 {
	return abc.GetPrice(symbol)
}

func (abc *ArbitrageKucoinClient) Ping() bool {
	return abc.Ping()
}
