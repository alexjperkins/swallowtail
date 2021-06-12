package domain
// BinanceEvent encapsulates the JSON response from Binance streams type BinanceEvent struct {
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

// BinanceTradeEvent is a trade event on binance
type BinanceTradeEvent struct {
	EventType          string `json:"e"`
	EventTime          int    `json:"E"`
	Symbol             string `json:"s"`
	TradeID            int    `json:"t"`
	Price              string `json:"p"`
	Quantity           string `json:"q"`
	BuyerOrderID       int    `json:"b"`
	SellerOrderID      int    `json:"a"`
	TradeTime          int    `json:"T"`
	IsBuyerMarketMaker bool   `json:"m"`
}

// BinanceKlineEvent is a Kline/Candlestick Event.
type BinanceKlineEvent struct {
	EventType string                `json:"e"`
	EventTime int                   `json:"E"`
	Symbol    string                `json:"s"`
	Data      BinanceKlineEventData `json:"k"`
}

type BinanceKlineEventData struct {
	KlineStartTime  int    `json:"t"`
	KlineCloseTime  int    `json:"T"`
	Symbol          string `json:"s"`
	Interval        string `json:"i"`
	FirstTradeID    int    `json:"f"`
	LastTradeID     int    `json:"L"`
	OpenPrice       string `json:"o"`
	ClosePrice      string `json:"c"`
	LowPrice        string `json:"l"`
	BaseAssetVolume string `json:"v"`
	NumberOfTrade   int    `json:"n"`

	IsKlineClosed            bool   `json:"x"`
	QuoteAssetVolume         string `json:"q"`
	TakerBuyBaseAssetVolume  string `json:"V"`
	TakerBuyQuoteAssetVolume string `json:"Q"`
}

// BinaceStreamPong
type BinanceStreamPong struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}
