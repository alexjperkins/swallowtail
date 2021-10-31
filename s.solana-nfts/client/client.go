package client

import (
	"context"
	"swallowtail/libraries/gerrors"
)

var (
	client SolanaNFTVendorClient
)

// SolanaNFTVendorClient ...
type SolanaNFTVendorClient interface {
	Ping(ctx context.Context) error
}

// Init initializes are required vendor clients
func Init(ctx context.Context) error {
	c := vendorClient{}

	if err := c.Ping(ctx); err != nil {
		return gerrors.Augment(err, "failed_to_establish_connection_to_all_vendor_clients", nil)
	}

	return nil
}
