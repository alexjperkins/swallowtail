package arbitrage

import (
	"swallowtail/s.twitter/clients"
)

var (
	ArbitrageCoinbaseClientID = "arbitrage-coinbase-client"
)

func init() {
	register(ArbitrageCoinbaseClientID,
		&ArbitrageCoinbaseClient{
			// TODO: thread default timeout
			i: &clients.CoinbaseClient{},
		},
	)
}

type ArbitrageCoinbaseClient struct {
	i *clients.CoinbaseClient
}

func (abc *ArbitrageCoinbaseClient) GetPrice(symbol string) float64 {
	return abc.GetPrice(symbol)
}

func (abc *ArbitrageCoinbaseClient) Ping() bool {
	return abc.Ping()
}
