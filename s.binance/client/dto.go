package client

import "time"

// Credentials ...
type Credentials struct {
	APIKey    string
	SecretKey string
}

// AsHeaders helper function that converts the underlying credentials struct to headers format for the Binance API.
func (c *Credentials) AsHeaders() map[string]string {
	return map[string]string{
		"X-MBX-APIKEY": c.APIKey,
		"Content-Type": "application/x-www-form-urlencoded",
	}
}

// BinanceAssetItem defines the asset item.
type BinanceAssetItem struct {
	// e.g BTCUSDT
	Symbol string `json:"symbol"`
	// e.g ETH
	BaseAsset string `json:"baseAsset"`
	// e.g USDT
	QuoteAsset        string `json:"quoteAsset"`
	WithMarginTrading bool   `json:"isMarginTradingAllowed"`
	WithSpotTrading   bool   `json:"isSpotTradingAllowed"`
}

// GetLatestPriceRequest ...
type GetLatestPriceRequest struct {
	Symbol string `json:"symbols"`
}

// GetLatestPriceResponse ...
type GetLatestPriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Time   int    `json:"time"`
}

// ListAllAssetPairsRequest is the data transfer object for listing all asset pairs on Binance.
type ListAllAssetPairsRequest struct{}

// ListAllAssetPairsResponse the response definition for `ListAllAssetPairs`
type ListAllAssetPairsResponse struct {
	Symbols []*BinanceAssetItem `json:"symbols"`
}

type ExecuteSpotTradeRequest struct{}
type ExecuteSpotTradeResponse struct{}

// PingRequest data transfer object for ping request.
type PingRequest struct{}

// PingResponse data transfer object for ping response.
type PingResponse struct{}

type ReadSpotAccountRequest struct{}

type ReadSpotAccountResponse struct{}

/// -- Perpetual Futures --- ///

// ReadPerpetualFuturesAccountRequest ...
type ReadPerpetualFuturesAccountRequest struct{}

// PerpetualFuturesAccountBalance
type PerpetualFuturesAccountBalance struct {
	AccountAlias       string `json:"accountAlias"`
	Asset              string `json:"asset"`
	Balance            string `json:"balance"`
	CrossWalletBalance string `json:"crossWalletBalance"`
	CrossUnPNL         string `json:"crossUnPnl"`
	AvailableBalance   string `json:"availableBalance"`
	MaxWithdrawAmount  string `json:"maxWithdrawAmount"`
	MarginAvailable    bool   `json:"marginAvailable"`
	LastUpdated        int    `json:"updateTime"`
}

// ReadPerpetualFuturesAccountResponse an array of balances; identical to the Binance exchange API definition.
type ReadPerpetualFuturesAccountResponse []*PerpetualFuturesAccountBalance

// ExecutePerpetualFuturesTradeRequest...
// https://binance-docs.github.io/apidocs/futures/en/#place-multiple-orders-trade
type ExecutePerpetualFuturesTradeRequest struct {
	// Required
	Symbol string `json:"string"`
	Side   string `json:"side"`
	Type   string `json:"type"`
	// ---
	PositionSide     string  `json:"positionSide"` // "BOTH", "LONG", "SHORT"
	TimeInForce      string  `json:"timeInForce"`
	Quantity         string  `json:"quantity"`
	ReduceOnly       string  `json:"reduceOnly"` // "true" or "false"
	Price            string  `json:"price"`
	NewClientOrderID string  `json:"newClientOrderId"`
	StopPrice        string  `json:"stopPrice"`
	ClosePosition    string  `json:"closePosition"`
	ActivationPrice  float64 `json:"activationPrice"`
	CallbackRate     float64 `json:"callbackRate"` // Used with trailing stop
	WorkingType      string  `json:"workingType"`  // "MARK_PRICE" or "CONTRACT_PRICE"
	PriceProtect     string  `json:"priceProtect"`
	NewOrderRespType string  `json:"newOrderRespType"` // "ACK" or "RESULT"
}

// ExecutePerpetualFuturesTradeResponse ...
type ExecutePerpetualFuturesTradeResponse struct {
	OrderID            int `json:"orderId"`
	ExecutionTimestamp int `json:"updateTime"`
}

// VerifyCredentialsRequest ...
type VerifyCredentialsRequest struct {
	Credentials *Credentials
}

// VerifyCredentialsResponse ...
type VerifyCredentialsResponse struct {
	IPRestrict                     bool `json:"ipRestrict"`
	CreateTime                     int  `json:"createTime"`
	EnableWithdrawals              bool `json:"enableWithdrawals"`
	EnableInternalTransfer         bool `json:"enableInternalTransfer"`
	PermitsUniversalTransfer       bool `json:"permitsUniversalTransfer"`
	EnableVanillaOptions           bool `json:"enableVanillaOptions"`
	EnableReading                  bool `json:"enableReading"`
	EnableFutures                  bool `json:"enableFutures"`
	EnableMargin                   bool `json:"enableMargin"`
	EnableSpotAndMarginTrading     bool `json:"enableSpotAndMarginTrading"`
	TradingAuthorityExpirationTime int  `json:"tradingAuthorityExpirationTime"`
}

// GetStatusRequest ...
type GetStatusRequest struct{}

// GetStatusResponse ...
type GetStatusResponse struct {
	AssumedClockDrift time.Duration
	ServerTime        int `json:"serverTime"`
	ServerLatency     time.Duration
}

// GetFuturesExchangeInfoRequest ...
type GetFuturesExchangeInfoRequest struct{}

// GetFuturesExchangeInfoResponse ...
type GetFuturesExchangeInfoResponse struct {
	RateLimits []struct {
		Interval      string `json:"interval"`
		IntervalNum   int    `json:"intervalNum"`
		Limit         int    `json:"limit"`
		RateLimitType string `json:"rate_limit_type"`
	} `json:"rate_limits"`
	ServerTime int `json:"serverTime"`
	Assets     []struct {
		Asset           string `json:"asset"`
		MarginAvailable bool   `json:"margainAvailable"`
	} `json:"assets"`
	Symbols []struct {
		Symbol            string `json:"symbol"`
		Pair              string `json:"pair"`
		ContractType      string `json:"contractType"`
		BaseAsset         string `json:"baseAsset"`
		QuoteAsset        string `json:"quoteAsset"`
		MarginAsset       string `json:"marginAsset"`
		QuantityPrecision int    `json:"quantityPrecision"`
		PricePrecision    int    `json:"pricePrecision"`
	}
}
