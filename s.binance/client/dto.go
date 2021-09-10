package client

// Credentials
type Credentials struct {
	APIKey    string
	SecretKey string
}

// AsHeaders helper function that converts the underlying credentials struct to headers format for the Binance API.
func (c *Credentials) AsHeaders() map[string]string {
	return map[string]string{
		"X-MBX-APIKEY": c.APIKey,
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

type ReadPerptualFuturesAccountRequest struct{}

type ReadPerptualFuturesAccountResponse struct{}

type VerifyCredentialsRequest struct {
	Credentials *Credentials
}

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
