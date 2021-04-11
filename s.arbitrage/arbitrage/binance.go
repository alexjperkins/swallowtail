package arbitrage

import (
	"swallowtail/s.twitter/clients"
)

var (
	ArbitrageBinanceClientID = "arbitrage-binance-client"
)

func init() {
	register(ArbitrageBinanceClientID,
		&ArbitrageBinanceClient{
			// TODO: thread default timeout
			i: clients.NewBinanceClient(),
		},
	)
}

type ArbitrageBinanceClient struct {
	i *clients.BinanceClient
}

func (abc *ArbitrageBinanceClient) GetPrice(symbol string) float64 {
	return abc.GetPrice(symbol)
}

func (abc *ArbitrageBinanceClient) Ping() bool {
	return abc.Ping()
}
