package client

import (
	"fmt"
	"time"
)

// PaginationFilter provides a pagination filter for all requests that require it.
type PaginationFilter struct {
	// Both Start & End are of second granularity.
	Start int64
	End   int64
}

// ToQueryString ...
func (p *PaginationFilter) ToQueryString() string {
	return fmt.Sprintf("start_time=%d&end_time=%d", p.Start, p.End)
}

// ExecuteOrderRequest ...
type ExecuteOrderRequest struct {
	ClientID         string  `json:"clientId,omitempty"`
	Market           string  `json:"market"`
	Side             string  `json:"side"`
	Type             string  `json:"type"`
	Size             float64 `json:"size"`
	Price            float64 `json:"price"`
	ReduceOnly       bool    `json:"reduceOnly,omitempty"`
	IOC              bool    `json:"ioc,omitempty"`
	PostOnly         bool    `json:"postOnly,omitempty"`
	RetryUntilFilled bool    `json:"retryUntilFilled,omitempty"`
	TriggerPrice     float64 `json:"triggerPrice,omitempty"`
	OrderPrice       float64 `json:"orderPrice,omitempty"`
	TrailValue       float64 `json:"trailValue,omitempty"`
}

// ExecuteOrderResponse ...
type ExecuteOrderResponse struct {
	Success bool `json:"success"`
	Result  struct {
		CreatedAt     time.Time `json:"createdAt,omitempty"`
		FilledSize    float64   `json:"filledSize,omitempty"`
		Future        string    `json:"future,omitempty"`
		ID            int       `json:"id,omitempty"`
		Market        string    `json:"market,omitempty"`
		Price         float64   `json:"price,omitempty"`
		RemainingSize float64   `json:"remainingSize,omitempty"`
		Side          string    `json:"side,omitempty"`
		Quantity      float64   `json:"size,omitempty"`
		Status        string    `json:"status,omitempty"`
		Type          string    `json:"type,omitempty"`
		ReduceOnly    bool      `json:"reduceOnly,omitempty"`
		IOC           bool      `json:"ioc,omitempty"`
		PostOnly      bool      `json:"postOnly,omitempty"`
		ClientID      string    `json:"clientId,omitempty"`
	} `json:"result"`
}

// GetStatusRequest ...
type GetStatusRequest struct{}

// GetStatusResponse ...
type GetStatusResponse struct {
	Success bool `json:"success"`
	Result  []struct {
		P50Latency float64 `json:"p50"`
	} `json:"result"`
}

// ListAccountDepositsRequest ...
type ListAccountDepositsRequest struct{}

// ListAccountDepositsResponse ...
type ListAccountDepositsResponse struct {
	Success  bool             `json:"success"`
	Deposits []*DepositRecord `json:"result"`
}

// DepositRecord ...
type DepositRecord struct {
	ID            int64     `json:"id"`
	Coin          string    `json:"coin"`
	TXID          string    `json:"txid"`
	Size          float64   `json:"size"`
	Fee           float64   `json:"fee"`
	Status        string    `json:"status"`
	Time          time.Time `json:"time"`
	SentTime      time.Time `json:"sentTime"`
	ConfirmedTime time.Time `json:"confirmedTime"`
	Confirmations int64     `json:"confirmations"`
}

// VerifyCredentialsRequest ...
type VerifyCredentialsRequest struct{}

// VerifyCredentialsResponse ...
type VerifyCredentialsResponse struct {
	Success bool `json:"success"`
}

// GetFundingRateRequest ...
type GetFundingRateRequest struct {
	Instrument string `json:"future"`
	StartTime  int    `json:"start_time"`
	EndTime    int    `json:"end_time"`
	Limit      int    `json:"limit"`
}

// FundingRateInfo ...
type FundingRateInfo struct {
	Instrument string  `json:"future"`
	Rate       float64 `json:"rate"`
	Time       string  `json:"time"`
}

// GetFundingRateResponse ...
type GetFundingRateResponse struct {
	FundingRates []*FundingRateInfo `json:"result"`
	Success      bool               `json:"success"`
}

// ListInstrumentsRequest ...
// https://docs.ftx.com/?python#markets
type ListInstrumentsRequest struct{}

// Instrument ...
type Instrument struct {
	Symbol          string  `json:"name"`
	PostOnly        bool    `json:"postOnly"`
	Enabled         bool    `json:"enabled"`
	MininumTickSize float64 `json:"priceIncrement"`
	MininumQuantity float64 `json:"sizeIncrement"`
	Type            string  `json:"type"`
	Underlying      string  `json:"underlying"`
	BaseCurrency    string  `json:"baseCurrency,omitempty"`
	QuoteCurrency   string  `json:"quoteCurrency,omitempty"`
}

// ListInstrumentsResponse ...
type ListInstrumentsResponse struct {
	Instruments []*Instrument `json:"result"`
}

// ReadAccountInformationResponseResult ...
type ReadAccountInformationResponseResult struct {
	BackstopProvider             bool    `json:"backstopProvider,omitempty"`
	Collateral                   float64 `json:"collateral,omitempty"`
	FreeCollateral               float64 `json:"freeCollateral,omitempty"`
	InitialMarginRequirement     float64 `json:"initialMarginRequirement,omitempty"`
	Leverage                     float64 `json:"leverage,omitempty"`
	Liquidating                  bool    `json:"liquidating,omitempty"`
	MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement,omitempty"`
	MakerFee                     float64 `json:"makerFee,omitempty"`
	MarginFraction               float64 `json:"marginFraction,omitempty"`
	OpenMarginFraction           float64 `json:"openMarginFraction,omitempty"`
	TakerFee                     float64 `json:"takerFee,omitempty"`
	TotalAccountValue            float64 `json:"totalAccountValue,omitempty"`
	TotalPositionSize            float64 `json:"totalPositionSize,omitempty"`
	Username                     string  `json:"username,omitempty"`
}

// ReadAccountInformationResponse ...
type ReadAccountInformationResponse struct {
	Success bool                                  `json:"success"`
	Result  *ReadAccountInformationResponseResult `json:"result"`
}

// AccountBalance ...
type AccountBalance struct {
	Asset                  string  `json:"coin"`
	QuantityAvailable      float64 `json:"free"`
	SpotBorrow             float64 `json:"spotBorrow,omitempty"`
	Total                  float64 `json:"total"`
	USDValue               float64 `json:"usdValue,omitempty"`
	AvailableWithoutBorrow float64 `json:"availableWithoutBorrow"`
}

// ListAccountBalancesResponse ...
type ListAccountBalancesResponse struct {
	Success         bool              `json:"success"`
	AccountBalances []*AccountBalance `json:"result"`
}
