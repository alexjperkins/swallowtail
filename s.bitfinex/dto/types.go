package dto

// GetStatusRequest ...
type GetStatusRequest struct {
}

// GetStatusProxyResponse ...
type GetStatusProxyResponse [1]int

// Operative ...
func (p *GetStatusProxyResponse) Operative() int {
	return p[0]
}

// GetStatusResponse ...
type GetStatusResponse struct {
	// Defines if the Bitfinex platfrom is live.
	Operative int `json:"operative"`
	// Latency of the server.
	ServerLatency int `json:"-"`
}

// GetFundingRatesRequest ...
type GetFundingRatesRequest struct {
	Symbol string `json:"symbol"`
}

// GetFundingRatesProxyResponse ...
type GetFundingRatesProxyResponse [][24]interface{}

// CurrentFundingRate ...
func (p GetFundingRatesProxyResponse) CurrentFundingRate() float64 {
	if len(p) < 1 {
		return 0.0
	}

	f, ok := p[0][9].(float64)
	if !ok {
		return 0.0
	}

	return f
}

// FundingRateInfo ...
type FundingRateInfo struct {
	FundingRate float64 `json:"funding_rate"`
}

// GetFundingRatesResponse ...
type GetFundingRatesResponse struct {
	Symbol       string             `json:"symbol"`
	FundingRates []*FundingRateInfo `json:"funding_rates"`
}
