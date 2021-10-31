package client

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"

	"swallowtail/libraries/ratelimit"
	"swallowtail/libraries/transport"
)

// NewVendorClient ...
func NewVendorClient(ctx context.Context, http transport.HttpClient) *VendorClient {
	return &VendorClient{
		&solanartClient{rateLimiter: ratelimit.NewLinearBackpressureRateLimiter(ctx, time.Minute, 30), http: http},
		&magicEdenClient{rateLimiter: ratelimit.NewLinearBackpressureRateLimiter(ctx, time.Minute, 30), http: http},
	}
}

// VendorClient is an abstraction to group all different vendor clients
// the sum of behaviour must equal the interface defined.
type VendorClient struct {
	*solanartClient
	*magicEdenClient
}

// Ping ...
func (v *VendorClient) Ping(ctx context.Context) error {
	var mErr error
	if err := v.solanartClient.Ping(ctx); err != nil {
		mErr = multierror.Append(err)
	}

	if err := v.magicEdenClient.Ping(ctx); err != nil {
		mErr = multierror.Append(err)
	}

	return mErr
}
