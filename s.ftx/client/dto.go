package client

import (
	"fmt"
	"swallowtail/s.ftx/client/signer"
	"time"
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
	return signer.Sha256HMAC(c.SecretKey, signaturePayload)
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

// PingRequest ...
type PingRequest struct{}

// PingResponse ...
type PingResponse struct {
	Success bool `json:"success"`
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
