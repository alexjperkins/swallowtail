package client

import (
	"context"
	"net/url"
	"swallowtail/libraries/gerrors"
	"swallowtail/libraries/transport"
	"swallowtail/libraries/util"
	"time"

	"github.com/opentracing/opentracing-go"
)

var (
	defaultHostname           = "https://ftx.com"
	client                    FTXClient
	depositAccountCredentials *Credentials
)

// FTXClient defines the client side contract of the FTX REST API
type FTXClient interface {
	// Ping ...
	Ping(ctx context.Context) error

	// GetStatus ...
	GetStatus(ctx context.Context, req *GetStatusRequest) (*GetStatusResponse, error)

	// ListAccountDeposits ...
	ListAccountDeposits(ctx context.Context, req *ListAccountDepositsRequest, pagination *PaginationFilter) (*ListAccountDepositsResponse, error)

	// VerifyCredentials ...
	VerifyCredentials(ctx context.Context, req *VerifyCredentialsRequest, credentials *Credentials) (*VerifyCredentialsResponse, error)

	// GetFundingRate ...
	GetFundingRate(ctx context.Context, req *GetFundingRateRequest) (*GetFundingRateResponse, error)
}

// Init instantiates the FTX client singleton.
func Init(ctx context.Context) error {
	c := &ftxClient{
		http:     transport.NewHTTPClient(30*time.Second, nil),
		hostname: defaultHostname,
	}

	if err := c.Ping(ctx); err != nil {
		return gerrors.Augment(err, "failed_to_init_ftx_client", nil)
	}

	apiKey := util.SetEnv("FTX_DEPOSIT_ACCOUNT_API_KEY")
	secretKey := util.SetEnv("FTX_DEPOSIT_ACCOUNT_SECRET_KEY")
	subaccount := util.SetEnv("FTX_DEPOSIT_ACCOUNT_SUBACCOUNT")

	if apiKey == "" || secretKey == "" || subaccount == "" {
		return gerrors.FailedPrecondition("failed_to_init_ftx_client.deposit_account_credentials_not_set", nil)
	}

	depositAccountCredentials = &Credentials{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		Subaccount: url.PathEscape(subaccount),
	}

	if _, err := c.VerifyCredentials(ctx, &VerifyCredentialsRequest{}, depositAccountCredentials); err != nil {
		return gerrors.Augment(err, "failed_to_init_ftx_client.failed_to_verify_deposit_account_credentials", map[string]string{
			"subaccount": depositAccountCredentials.Subaccount,
		})
	}

	client = c

	return nil
}

// ListAccountDeposits ...
func ListAccountDeposits(ctx context.Context, req *ListAccountDepositsRequest, pagination *PaginationFilter) (*ListAccountDepositsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "List account deposits on FTX")
	defer span.Finish()
	return client.ListAccountDeposits(ctx, req, pagination)
}

// VerifyCredentials ...
func VerifyCredentials(ctx context.Context, req *VerifyCredentialsRequest, credentials *Credentials) (*VerifyCredentialsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Verify FTX credentials deposits")
	defer span.Finish()
	return client.VerifyCredentials(ctx, req, credentials)

}

// GetFundingRate ...
func GetFundingRate(ctx context.Context, req *GetFundingRateRequest) (*GetFundingRateResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get funding rate FTX")
	defer span.Finish()
	return client.GetFundingRate(ctx, req)
}

// GetStatus ...
func GetStatus(ctx context.Context, req *GetStatusRequest) (*GetStatusResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get status: FTX")
	defer span.Finish()
	return client.GetStatus(ctx, req)
}
