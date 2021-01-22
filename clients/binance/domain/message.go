package domain

// BinanceMsg encapsulates the JSON response from Binance streams
type BinanceMsg struct {
	EventType                   string `json:"e"`
	EventTime                   int    `json:"E"`
	Symbol                      string `json:"s"`
	PriceDelta                  string `json:"p"`
	PriceDeltaPercentage        string `json:"P"`
	WeightedAveragePrice        string `json:"w"`
	StartPrice                  string `json:"x"`
	LastPrice                   string `json:"c"`
	LastQuantity                string `json:"Q"`
	BestBidPrice                string `json:"b"`
	BestBidQuantity             string `json:"B"`
	BestAskPrice                string `json:"a"`
	BestAskQuantity             string `json:"A"`
	OpenPrice                   string `json:"o"`
	HighPrice                   string `json:"h"`
	LowPrice                    string `json:"l"`
	TotalTradedBaseAssetVolume  string `json:"v"`
	TotalTradedQuoteAssetVolume string `json:"q"`
	OpenTime                    int    `json:"O"`
	CloseTime                   int    `json:"C"`
	FirstTradeID                int    `json:"F"`
	LastTradeID                 int    `json:"L"`
	TotalNumberOfTrades         int    `json:"n"`
}
