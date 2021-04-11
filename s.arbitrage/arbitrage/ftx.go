package arbitrage

import (
	"swallowtail/s.twitter/clients"
)

var (
	ArbitrageFTXClientID = "arbitrage-ftx-client"
)

func init() {
	register(ArbitrageFTXClientID,
		&ArbitrageFTXClient{
			// TODO: thread default timeout
			i: &clients.FTXClient{},
		},
	)
}

type ArbitrageFTXClient struct {
	i *clients.FTXClient
}

func (abc *ArbitrageFTXClient) GetPrice(symbol string) float64 {
	return abc.GetPrice(symbol)
}

func (abc *ArbitrageFTXClient) Ping() bool {
	return abc.Ping()
}
