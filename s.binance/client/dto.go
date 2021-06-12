package client

// BinanceAssetItem defines the asset item.
type BinanceAssetItem struct {
	Symbol            string `json:"symbol"`
	BaseAsset         string `json:"baseAsset"`
	WithMarginTrading bool   `json:"isMarginTradingAllowed"`
	WithSpotTrading   bool   `json:"isSpotTradingAllowed"`
}

// ListAllAssetPairsRequest is the data transfer object for listing all asset pairs on Binance.
type ListAllAssetPairsRequest struct{}

// ListAllAssetPairsResponse the response definition for `ListAllAssetPairs`
type ListAllAssetPairsResponse struct {
	Symbols []*BinanceAssetItem `json:"symbols"`
}

// PingRequest data transfer object for ping request.
type PingRequest struct{}

// PingResponse data transfer object for ping response.
type PingResponse struct{}
