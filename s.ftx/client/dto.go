package client

import (
	"fmt"
	"time"

	"swallowtail/s.ftx/client/auth"
)

// Credentials holds the credentials required for FTX.
type Credentials struct {
	APIKey     string
	SecretKey  string
	Subaccount string
}

// SignRequest signs a given request with the request body.
func (c *Credentials) SignRequest(method, endpoint, timestamp string, body []byte) (string, error) {
	signaturePayload := fmt.Sprintf("%s%s%s%s", timestamp, method, endpoint, body)
	return auth.Sha256HMAC(c.SecretKey, signaturePayload)
}

// AsHeaders converts the credentials struct into the headers required to verify the user.
// It uses the request body & the timestamp to sign the request.
func (c *Credentials) AsHeaders(signature, timestamp string) map[string]string {
	m := map[string]string{
		"Content-Type": "application/json",
		"FTX-KEY":      c.APIKey,
		"FTX-SIGN":     signature,
		"FTX-TS":       timestamp,
	}

	if c.Subaccount != "" {
		m["FTX-SUBACCOUNT"] = c.Subaccount
	}

	return m
}

// SubaccountAsHeaders returns only the subaccount as headers; if it is not null.
func (c *Credentials) SubaccountAsHeaders() map[string]string {
	m := map[string]string{
		"Content-Type": "application/json",
	}
	if c.Subaccount != "" {
		m["FTX-SUBACCOUNT"] = c.Subaccount
	}
	return m
}

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
	Symbol            string `json:"symbol"`
	Side              string `json:"side"`
	Price             string `json:"price"`
	Type              string `json:"type"`
	Quantity          string `json:"quantity"`
	ReduceOnly        bool   `json:"reduce_only"`
	IOC               bool   `json:"ioc"`
	PostOnly          bool   `json:"post_only"`
	RejectOnPriceBand bool   `json:"reject_on_price_band"`
	RetryUntilFilled  bool   `json:"retry_until_filled"`
	TriggerPrice      string `json:"trigger_price"`
	OrderPrice        string `json:"order_price"`
	TrailValue        string `json:"trail_value"`
}

// ExecuteOrderResponse ...
type ExecuteOrderResponse struct {
	Success bool `json:"success"`
	Result  struct {
		CreatedAt     time.Time `json:"createdAt"`
		FilledSize    float64   `json:"filledSize"`
		Future        string    `json:"future"`
		ID            int       `json:"id"`
		Market        string    `json:"market"`
		Price         float64   `json:"price"`
		RemainingSize float64   `json:"remaining_size"`
		Side          string    `json:"side"`
		Quantity      float64   `json:"size"`
		Status        string    `json:"status"`
		Type          string    `json:"type"`
		ReduceOnly    bool      `json:"reduce_only"`
		IOC           bool      `json:"ioc"`
		PostOnly      bool      `json:"post_only"`
		ClientID      string    `json:"client_id"`
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

// Deposit ...
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
type ListInstrumentsRequest struct {
}

// ListFuturesInstrumentsResponse ...
// https://docs.ftx.com/?python#futures
type ListFuturesInstrumentsResponse struct {
}

// ListMarketsInstrumentsResponse ...
// https://docs.ftx.com/?python#markets
type ListMarketsInstrumentsResponse struct {
}

// Instrument ...
type Instrument struct {
	Symbol string `json:"symbol"`
}

// ListInstrumentsResponse ...
type ListInstrumentsResponse struct {
	Instruments []*Instrument
}
